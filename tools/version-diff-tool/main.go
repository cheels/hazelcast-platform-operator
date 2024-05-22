package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

type CRD struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Names struct {
			Kind string `yaml:"kind"`
		} `yaml:"names"`
		Versions []struct {
			Schema struct {
				OpenAPIV3Schema struct {
					Properties map[string]struct {
						Description string                 `yaml:"description,omitempty"`
						Properties  map[string]interface{} `yaml:"properties"`
						Required    []string               `yaml:"required,omitempty"`
						Type        string                 `yaml:"type"`
					} `yaml:"properties"`
				} `yaml:"openAPIV3Schema"`
			} `yaml:"schema"`
		} `yaml:"versions"`
	} `yaml:"spec"`
}

func createOpenAPISpec(crds []CRD) map[string]interface{} {
	paths := make(map[string]interface{})
	for _, crd := range crds {
		path := fmt.Sprintf("/%s.%s", strings.ToLower(crd.Spec.Names.Kind), "hazelcast.com")
		schema := map[string]interface{}{
			"description": crd.Metadata.Name,
			"type":        "object",
			"required":    []string{"spec"},
			"properties": map[string]interface{}{
				"apiVersion": map[string]interface{}{"type": "string"},
				"kind":       map[string]interface{}{"type": "string"},
				"metadata":   map[string]interface{}{"type": "object"},
				"spec":       crd.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"],
				"status":     crd.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["status"],
			},
		}

		paths[path] = map[string]interface{}{
			"post": map[string]interface{}{
				"requestBody": map[string]interface{}{
					"required": true,
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": schema,
						},
					},
				},
			},
		}
	}
	return map[string]interface{}{"paths": paths}
}

func addAndUpdateRepo(settings *cli.EnvSettings, repoName, repoURL string) error {
	repoEntry := &repo.Entry{
		Name: repoName,
		URL:  repoURL,
	}
	chartRepo, err := repo.NewChartRepository(repoEntry, getter.All(settings))
	if err != nil {
		return fmt.Errorf("failed to create chart repository: %v", err)
	}
	if _, err := chartRepo.DownloadIndexFile(); err != nil {
		return fmt.Errorf("failed to download index file: %v", err)
	}
	return nil
}

func generateCRDFile(version string) (string, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), nil); err != nil {
		return "", fmt.Errorf("failed to initialize Helm action configuration: %v", err)
	}

	repoName := "hazelcast"
	repoURL := "https://hazelcast-charts.s3.amazonaws.com"
	crdChartName := "hazelcast-platform-operator-crds"

	if err := addAndUpdateRepo(settings, repoName, repoURL); err != nil {
		return "", err
	}

	client := action.NewInstall(actionConfig)
	client.DryRun = true
	client.ReleaseName = repoName
	client.Replace = true
	client.ClientOnly = true
	client.IncludeCRDs = true
	client.Version = version

	chartName := fmt.Sprintf("%s/%s", repoName, crdChartName)
	chartPath, err := client.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		return "", fmt.Errorf("failed to locate chart: %v", err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return "", fmt.Errorf("failed to load chart: %v", err)
	}

	release, err := client.Run(chart, nil)
	if err != nil {
		return "", fmt.Errorf("failed to render chart: %v", err)
	}

	outputFile := fmt.Sprintf("%s.yaml", version)
	if err := os.WriteFile(outputFile, []byte(release.Manifest), 0644); err != nil {
		return "", fmt.Errorf("failed to write CRD file to %s: %v", outputFile, err)
	}
	return outputFile, nil
}

func extractCRDs(inputFile string) ([]CRD, error) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var crds []CRD
	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	for {
		var crd CRD
		if err := decoder.Decode(&crd); err != nil {
			break
		}
		if crd.Kind == "CustomResourceDefinition" {
			crds = append(crds, crd)
		}
	}
	return crds, nil
}

func writeOpenAPISpec(outputFile string, openAPISpec map[string]interface{}) error {
	data, err := yaml.Marshal(&openAPISpec)
	if err != nil {
		return fmt.Errorf("failed to marshal OpenAPI spec: %v", err)
	}
	return os.WriteFile(outputFile, data, 0644)
}

func filterOutput(output string) string {
	warningMsgPattern := regexp.MustCompile(`This is a warning because.*?change in specification\.`)
	apiPattern := regexp.MustCompile(`in API`)
	postPattern := regexp.MustCompile(`POST`)
	output = warningMsgPattern.ReplaceAllString(output, "")
	output = apiPattern.ReplaceAllString(output, "")
	output = postPattern.ReplaceAllString(output, "in")
	output = strings.Replace(output, "/", "", 1)
	output = strings.Replace(output, "/", ".", 2)
	return output
}

func generateAndExtractCRDs(version string, outputFile string, wg *sync.WaitGroup, specInfo **load.SpecInfo, apiLoader *openapi3.Loader) {
	defer wg.Done()
	rawCrdFile, err := generateCRDFile(version)
	if err != nil {
		log.Fatalf("failed to generate CRD file for version %s: %v", version, err)
	}

	crds, err := extractCRDs(rawCrdFile)
	if err != nil {
		log.Fatalf("failed to extract CRDs from file for version %s: %v", version, err)
	}

	openAPISpec := createOpenAPISpec(crds)
	if err := writeOpenAPISpec(outputFile, openAPISpec); err != nil {
		log.Fatalf("failed to write OpenAPI spec for version %s: %v", version, err)
	}

	*specInfo, err = load.NewSpecInfo(apiLoader, load.NewSource(outputFile))
	if err != nil {
		log.Fatalf("failed to load spec info for version %s: %v", version, err)
	}
}

func main() {
	base := flag.String("base", "", "Version of the first CRD to compare")
	revision := flag.String("revision", "", "Version of the second CRD to compare")
	flag.Parse()

	if *base == "" || *revision == "" {
		log.Fatal("both versions (-base, -revision) must be provided")
	}

	apiLoader := openapi3.NewLoader()
	apiLoader.IsExternalRefsAllowed = true
	var wg sync.WaitGroup
	var baseSpec, revisionSpec *load.SpecInfo
	baseFile := fmt.Sprintf("%s.yaml", *base)
	revisionFile := fmt.Sprintf("%s.yaml", *revision)

	wg.Add(2)

	go generateAndExtractCRDs(*base, baseFile, &wg, &baseSpec, apiLoader)
	go generateAndExtractCRDs(*revision, revisionFile, &wg, &revisionSpec, apiLoader)

	wg.Wait()

	diffRes, operationsSources, err := diff.GetPathsDiff(diff.NewConfig(),
		[]*load.SpecInfo{baseSpec},
		[]*load.SpecInfo{revisionSpec},
	)
	if err != nil {
		log.Fatalf("diff failed with %v", err)
	}

	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), diffRes, operationsSources)
	if len(errs) > 0 {
		localization := checker.NewDefaultLocalizer()
		count := errs.GetLevelCount()
		fmt.Printf("Comparing CRD files: %s and %s\n", *base, *revision)
		result := fmt.Sprintf(localization("total-errors", len(errs), count[checker.ERR], "error", count[checker.WARN], "warning"))
		for _, bc := range errs {
			output := bc.MultiLineError(localization, checker.ColorAlways)
			filteredOutput := filterOutput(output)
			result += fmt.Sprintf("\n%s\n", filteredOutput)
		}
		fmt.Printf("\n%s", result)
	}
}
