package compile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kalo-build/plugin-openapi-morphe/pkg/compile/cfg"
)

func OpenAPIToMorphe(config cfg.OpenAPIToMorpheConfig) error {
	config.Resolve()
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	specPath := filepath.Join(config.InputDir, config.SpecFileName)
	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		return fmt.Errorf("failed to read spec file '%s': %w", specPath, err)
	}

	spec, err := ParseOpenAPISpec(specBytes)
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	enumResult := BuildEnums(spec, config.SchemasToSkip)

	structures := BuildStructures(spec, enumResult, config.SchemasToSkip)

	if len(structures) > 0 {
		if err := WriteStructures(config.OutputPath, structures); err != nil {
			return fmt.Errorf("write structures: %w", err)
		}
	}

	if len(enumResult.Enums) > 0 {
		if err := WriteEnums(config.OutputPath, enumResult.Enums); err != nil {
			return fmt.Errorf("write enums: %w", err)
		}
	}

	return nil
}
