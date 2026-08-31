package compile

import morpheyaml "github.com/kalo-build/morphe-go/pkg/yaml"

func OpenAPITypeToMorpheFieldType(types []string, format string) morpheyaml.StructureFieldType {
	if hasType(types, "string") {
		return stringFormatToMorpheType(format)
	}
	if hasType(types, "integer") {
		return morpheyaml.StructureFieldTypeInteger
	}
	if hasType(types, "number") {
		return morpheyaml.StructureFieldTypeFloat
	}
	if hasType(types, "boolean") {
		return morpheyaml.StructureFieldTypeBoolean
	}
	return morpheyaml.StructureFieldTypeString
}

func stringFormatToMorpheType(format string) morpheyaml.StructureFieldType {
	switch format {
	case "uuid":
		return morpheyaml.StructureFieldTypeUUID
	case "date-time":
		return morpheyaml.StructureFieldTypeTime
	case "date":
		return morpheyaml.StructureFieldTypeDate
	case "password":
		return morpheyaml.StructureFieldTypeProtected
	default:
		return morpheyaml.StructureFieldTypeString
	}
}
