package compile

import (
	"fmt"

	morpheyaml "github.com/kalo-build/morphe-go/pkg/yaml"
)

// BuildStructures converts parsed OpenAPI schemas into Morphe structures.
// Only schemas with properties (object types) are converted; schemas that
// are pure enums (no properties) are skipped.
func BuildStructures(spec *ParsedSpec, enumRefs *BuildEnumsResult, skip map[string]bool) map[string]morpheyaml.Structure {
	structures := make(map[string]morpheyaml.Structure)

	for _, schema := range spec.Schemas {
		if skip[schema.Name] {
			continue
		}
		if len(schema.Properties) == 0 {
			continue
		}

		structure := buildStructure(schema, enumRefs)
		structures[structure.Name] = structure
	}

	return structures
}

func buildStructure(schema ParsedSchema, enumRefs *BuildEnumsResult) morpheyaml.Structure {
	fields := make(map[string]morpheyaml.StructureField, len(schema.Properties))

	for _, prop := range schema.Properties {
		fieldName := PropertyNameToFieldName(prop.Name)
		field := buildStructureField(schema.Name, prop, enumRefs)
		fields[fieldName] = field
	}

	return morpheyaml.Structure{
		Name:   schema.Name,
		Fields: fields,
	}
}

func buildStructureField(schemaName string, prop ParsedProperty, enumRefs *BuildEnumsResult) morpheyaml.StructureField {
	fieldType := resolveFieldType(schemaName, prop, enumRefs)

	var attrs []string
	if !prop.Required {
		attrs = append(attrs, "optional")
	}

	return morpheyaml.StructureField{
		Type:       fieldType,
		Attributes: attrs,
	}
}

func resolveFieldType(schemaName string, prop ParsedProperty, enumRefs *BuildEnumsResult) morpheyaml.StructureFieldType {
	enumKey := fmt.Sprintf("%s.%s", schemaName, prop.Name)
	if enumName, ok := enumRefs.FieldRef[enumKey]; ok {
		return morpheyaml.StructureFieldType(enumName)
	}

	if prop.Ref != "" {
		return morpheyaml.StructureFieldType(RefToSchemaName(prop.Ref))
	}

	if prop.IsArray {
		if prop.ItemRef != "" {
			return morpheyaml.StructureFieldType(RefToSchemaName(prop.ItemRef))
		}
		return OpenAPITypeToMorpheFieldType(prop.ItemTypes, prop.ItemFormat)
	}

	return OpenAPITypeToMorpheFieldType(prop.Types, prop.Format)
}
