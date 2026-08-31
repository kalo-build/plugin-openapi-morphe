package compile

import (
	"fmt"
	"sort"

	morpheyaml "github.com/kalo-build/morphe-go/pkg/yaml"
)

// BuildEnumsResult holds extracted enums and a mapping from
// (schemaName, propertyName) to the generated enum name so that
// structure fields can reference them.
type BuildEnumsResult struct {
	Enums    map[string]morpheyaml.Enum
	FieldRef map[string]string // key: "SchemaName.PropertyName" → enum name
}

// BuildEnums scans all parsed schemas for inline string enums and
// extracts them into standalone Morphe enum definitions.
func BuildEnums(spec *ParsedSpec, skip map[string]bool) *BuildEnumsResult {
	result := &BuildEnumsResult{
		Enums:    make(map[string]morpheyaml.Enum),
		FieldRef: make(map[string]string),
	}

	for _, schema := range spec.Schemas {
		if skip[schema.Name] {
			continue
		}

		if len(schema.EnumValues) > 0 {
			buildTopLevelEnum(result, schema)
			continue
		}

		for _, prop := range schema.Properties {
			if len(prop.EnumValues) > 0 {
				buildInlineEnum(result, schema.Name, prop)
			}
		}
	}

	return result
}

func buildTopLevelEnum(result *BuildEnumsResult, schema ParsedSchema) {
	enumName := schema.Name
	entries := buildEnumEntries(schema.EnumValues)
	result.Enums[enumName] = morpheyaml.Enum{
		Name:    enumName,
		Type:    morpheyaml.EnumTypeString,
		Entries: entries,
	}
}

func buildInlineEnum(result *BuildEnumsResult, schemaName string, prop ParsedProperty) {
	enumName := schemaName + PropertyNameToFieldName(prop.Name)
	entries := buildEnumEntries(prop.EnumValues)

	result.Enums[enumName] = morpheyaml.Enum{
		Name:    enumName,
		Type:    morpheyaml.EnumTypeString,
		Entries: entries,
	}

	key := fmt.Sprintf("%s.%s", schemaName, prop.Name)
	result.FieldRef[key] = enumName
}

func buildEnumEntries(values []string) map[string]any {
	entries := make(map[string]any, len(values))
	sorted := make([]string, len(values))
	copy(sorted, values)
	sort.Strings(sorted)
	for _, v := range sorted {
		entries[v] = v
	}
	return entries
}
