// main.go

package main

import (
	"flag"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gonvenience/bunt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/formatters"
	"github.com/tufin/oasdiff/load"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	white = bunt.GhostWhite
)
var gitHubFormatter = formatters.GitHubActionsFormatter{
	Localizer: checker.NewDefaultLocalizer(),
}

func colorText(format string, color colorful.Color, a ...interface{}) string {
	return bunt.Style(
		fmt.Sprintf(format, a...),
		bunt.EachLine(),
		bunt.Foreground(color),
	)
}
func underlineText(format string, a ...interface{}) string {
	return bunt.Style(
		fmt.Sprintf(format, a...),
		bunt.EachLine(),
		bunt.Underline(),
	)
}

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

func createOpenAPIBundle(crds []CRD) map[string]interface{} {
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

func generateCRDFile(version string) (string, error) {
	outputFile := fmt.Sprintf("%s.yaml", version)
	addRepoCmd := exec.Command("sh", "-c", "helm repo add hazelcast https://hazelcast-charts.s3.amazonaws.com && helm repo update")
	err := addRepoCmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to add and update hazelcast repo: %v", err)
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("helm template operator hazelcast/hazelcast-platform-operator-crds --version=%s > %s", version, outputFile))
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to generate CRD file for version %v", err)
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

func writeOpenAPIBundle(outputFile string, openAPISpec map[string]interface{}) error {
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
	crdNameRegex := regexp.MustCompile(`\b(\w+\.\w+\.\w+)\b`)

	output = warningMsgPattern.ReplaceAllString(output, "")
	output = apiPattern.ReplaceAllString(output, "")
	output = postPattern.ReplaceAllString(output, colorText("in", white))
	output = crdNameRegex.ReplaceAllString(output, fmt.Sprintf("%s", underlineText("$1")))
	output = strings.Replace(output, "/", "", 1)
	output = strings.Replace(output, "/", ".", 2)
	return output
}

func main() {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	base := flag.String("base", "", "Version of the first CRD to compare")
	revision := flag.String("revision", "", "Version of the second CRD to compare")
	//format := flag.Bool("f", false, "Apply filtering to the output and output only breaking changes")
	flag.Parse()

	if *base == "" || *revision == "" {
		log.Fatal("Both versions (-base, -revision) must be provided")
	}

	rawBaseCrd, err := generateCRDFile(*base)
	if err != nil {
		log.Fatalf("Failed to generate first CRD file: %v", err)
	}

	rawRevisionCrd, err := generateCRDFile(*revision)
	if err != nil {
		log.Fatalf("Failed to generate second CRD file: %v", err)
	}

	baseCrd, err := extractCRDs(rawBaseCrd)
	if err != nil {
		log.Fatalf("Failed to extract CRDs from first file: %v", err)
	}

	revisionCrd, err := extractCRDs(rawRevisionCrd)
	if err != nil {
		log.Fatalf("Failed to extract CRDs from second file: %v", err)
	}

	baseOpenAPISpec := createOpenAPIBundle(baseCrd)
	revisionOpenAPISpec := createOpenAPIBundle(revisionCrd)

	baseFile := fmt.Sprintf("%s.yaml", *base)
	if err := writeOpenAPIBundle(baseFile, baseOpenAPISpec); err != nil {
		log.Fatalf("Failed to write first OpenAPI bundle: %v", err)
	}

	revisionFile := fmt.Sprintf("%s.yaml", *revision)
	if err := writeOpenAPIBundle(revisionFile, revisionOpenAPISpec); err != nil {
		log.Fatalf("Failed to write second OpenAPI bundle: %v", err)
	}
	b, err := load.NewSpecInfo(loader, load.NewSource(baseFile))
	r, err := load.NewSpecInfo(loader, load.NewSource(revisionFile))
	diffRes, operationsSources, err := diff.GetPathsDiff(diff.NewConfig(),
		[]*load.SpecInfo{b},
		[]*load.SpecInfo{r},
	)
	if err != nil {
		log.Fatalf("diff failed with %v", os.Stderr)
		return
	}
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), diffRes, operationsSources)
	/*	errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.WARN, errs, "ignore-err.txt", checker.NewDefaultLocalizer())
		if err != nil {
			log.Fatalf("ignore errors failed with %v", os.Stderr)
			return
		}*/
	output, err := gitHubFormatter.RenderBreakingChanges(errs, formatters.NewRenderOpts())
	if len(errs) > 0 {
		localizer := checker.NewDefaultLocalizer()
		count := errs.GetLevelCount()
		fmt.Print(localizer("total-errors", len(errs), count[checker.ERR], "error", count[checker.WARN], "warning"))

		fmt.Printf("%s\n\n", output)
	}
}
