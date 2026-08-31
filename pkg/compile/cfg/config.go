package cfg

import "fmt"

type OpenAPIToMorpheConfig struct {
	InputDir     string
	OutputPath   string
	SpecFileName string
	SchemasToSkip map[string]bool
}

func (c *OpenAPIToMorpheConfig) Resolve() {
	if c.SpecFileName == "" {
		c.SpecFileName = "openapi.yaml"
	}
	if c.SchemasToSkip == nil {
		c.SchemasToSkip = make(map[string]bool)
	}
}

func (c *OpenAPIToMorpheConfig) Validate() error {
	if c.InputDir == "" {
		return fmt.Errorf("input directory is required")
	}
	if c.OutputPath == "" {
		return fmt.Errorf("output path is required")
	}
	return nil
}
