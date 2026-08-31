package compile_test

import (
	"testing"

	morpheyaml "github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-openapi-morphe/pkg/compile"
	"github.com/stretchr/testify/assert"
)

func TestOpenAPITypeToMorpheFieldType_String(t *testing.T) {
	result := compile.OpenAPITypeToMorpheFieldType([]string{"string"}, "")
	assert.Equal(t, morpheyaml.StructureFieldTypeString, result)
}

func TestOpenAPITypeToMorpheFieldType_UUID(t *testing.T) {
	result := compile.OpenAPITypeToMorpheFieldType([]string{"string"}, "uuid")
	assert.Equal(t, morpheyaml.StructureFieldTypeUUID, result)
}

func TestOpenAPITypeToMorpheFieldType_DateTime(t *testing.T) {
	result := compile.OpenAPITypeToMorpheFieldType([]string{"string"}, "date-time")
	assert.Equal(t, morpheyaml.StructureFieldTypeTime, result)
}

func TestOpenAPITypeToMorpheFieldType_Date(t *testing.T) {
	result := compile.OpenAPITypeToMorpheFieldType([]string{"string"}, "date")
	assert.Equal(t, morpheyaml.StructureFieldTypeDate, result)
}

func TestOpenAPITypeToMorpheFieldType_Password(t *testing.T) {
	result := compile.OpenAPITypeToMorpheFieldType([]string{"string"}, "password")
	assert.Equal(t, morpheyaml.StructureFieldTypeProtected, result)
}

func TestOpenAPITypeToMorpheFieldType_Integer(t *testing.T) {
	result := compile.OpenAPITypeToMorpheFieldType([]string{"integer"}, "")
	assert.Equal(t, morpheyaml.StructureFieldTypeInteger, result)
}

func TestOpenAPITypeToMorpheFieldType_IntegerInt32(t *testing.T) {
	result := compile.OpenAPITypeToMorpheFieldType([]string{"integer"}, "int32")
	assert.Equal(t, morpheyaml.StructureFieldTypeInteger, result)
}

func TestOpenAPITypeToMorpheFieldType_Number(t *testing.T) {
	result := compile.OpenAPITypeToMorpheFieldType([]string{"number"}, "")
	assert.Equal(t, morpheyaml.StructureFieldTypeFloat, result)
}

func TestOpenAPITypeToMorpheFieldType_NumberDouble(t *testing.T) {
	result := compile.OpenAPITypeToMorpheFieldType([]string{"number"}, "double")
	assert.Equal(t, morpheyaml.StructureFieldTypeFloat, result)
}

func TestOpenAPITypeToMorpheFieldType_Boolean(t *testing.T) {
	result := compile.OpenAPITypeToMorpheFieldType([]string{"boolean"}, "")
	assert.Equal(t, morpheyaml.StructureFieldTypeBoolean, result)
}

func TestOpenAPITypeToMorpheFieldType_UnknownDefaultsToString(t *testing.T) {
	result := compile.OpenAPITypeToMorpheFieldType([]string{"object"}, "")
	assert.Equal(t, morpheyaml.StructureFieldTypeString, result)
}

func TestOpenAPITypeToMorpheFieldType_EmptyDefaultsToString(t *testing.T) {
	result := compile.OpenAPITypeToMorpheFieldType(nil, "")
	assert.Equal(t, morpheyaml.StructureFieldTypeString, result)
}
