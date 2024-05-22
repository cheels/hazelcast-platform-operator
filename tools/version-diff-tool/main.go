package main

import (
	"flag"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
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

func generateCRDFile(version string) (string, error) {
	outputFile := fmt.Sprintf("%s.yaml", version)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("helm template hazelcast hazelcast/hazelcast-platform-operator-crds --version=%s > %s", version, outputFile))
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to generate CRD file for version %s: %v", version, err)
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
	crdNameRegex := regexp.MustCompile(`\b(\w+\.\w+\.\w+)\b`)

	output = warningMsgPattern.ReplaceAllString(output, "")
	output = apiPattern.ReplaceAllString(output, "")
	output = postPattern.ReplaceAllString(output, "in")
	output = crdNameRegex.ReplaceAllString(output, "$1")
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
		log.Fatal("both versions (-base, -revision) must be provided")
	}

	rawBaseCrd, err := generateCRDFile(*base)
	if err != nil {
		log.Fatalf("failed to generate first CRD file: %v", err)
	}

	rawRevisionCrd, err := generateCRDFile(*revision)
	if err != nil {
		log.Fatalf("failed to generate second CRD file: %v", err)
	}

	baseCrd, err := extractCRDs(rawBaseCrd)
	if err != nil {
		log.Fatalf("failed to extract CRDs from first file: %v", err)
	}

	revisionCrd, err := extractCRDs(rawRevisionCrd)
	if err != nil {
		log.Fatalf("failed to extract CRDs from second file: %v", err)
	}

	baseOpenAPISpec := createOpenAPISpec(baseCrd)
	revisionOpenAPISpec := createOpenAPISpec(revisionCrd)

	baseFile := fmt.Sprintf("%s.yaml", *base)
	if err := writeOpenAPISpec(baseFile, baseOpenAPISpec); err != nil {
		log.Fatalf("failed to write first OpenAPI spec: %v", err)
	}

	revisionFile := fmt.Sprintf("%s.yaml", *revision)
	if err := writeOpenAPISpec(revisionFile, revisionOpenAPISpec); err != nil {
		log.Fatalf("failed to write second OpenAPI : %v", err)
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
	if len(errs) > 0 || len(errs) == 0 {
		localizer := checker.NewDefaultLocalizer()
		count := errs.GetLevelCount()
		result := fmt.Sprintf(localizer("total-errors", len(errs), count[checker.ERR], "error", count[checker.WARN], "warning"))
		for _, bcerr := range errs {
			output := bcerr.MultiLineError(localizer, checker.ColorAlways)
			filteredOutput := filterOutput(output)
			result += fmt.Sprintf("\n%s\n", filteredOutput)
		}
		fmt.Printf("\n%s", result)
	}
}
