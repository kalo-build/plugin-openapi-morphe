package compile_test

import (
	"testing"

	morpheyaml "github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-openapi-morphe/pkg/compile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildEnums_InlineEnum(t *testing.T) {
	spec := &compile.ParsedSpec{
		Schemas: []compile.ParsedSchema{
			{
				Name: "Task",
				Properties: []compile.ParsedProperty{
					{
						Name:       "status",
						Types:      []string{"string"},
						EnumValues: []string{"open", "closed"},
						Required:   true,
					},
				},
			},
		},
	}

	result := compile.BuildEnums(spec, nil)

	require.Contains(t, result.Enums, "TaskStatus")
	enum := result.Enums["TaskStatus"]
	assert.Equal(t, "TaskStatus", enum.Name)
	assert.Equal(t, morpheyaml.EnumTypeString, enum.Type)
	assert.Equal(t, map[string]any{"closed": "closed", "open": "open"}, enum.Entries)

	assert.Equal(t, "TaskStatus", result.FieldRef["Task.status"])
}

func TestBuildEnums_TopLevelEnum(t *testing.T) {
	spec := &compile.ParsedSpec{
		Schemas: []compile.ParsedSchema{
			{
				Name:       "Color",
				EnumValues: []string{"red", "green", "blue"},
			},
		},
	}

	result := compile.BuildEnums(spec, nil)

	require.Contains(t, result.Enums, "Color")
	enum := result.Enums["Color"]
	assert.Equal(t, "Color", enum.Name)
	assert.Equal(t, morpheyaml.EnumTypeString, enum.Type)
	assert.Len(t, enum.Entries, 3)
}

func TestBuildEnums_SkippedSchema(t *testing.T) {
	spec := &compile.ParsedSpec{
		Schemas: []compile.ParsedSchema{
			{
				Name:       "Color",
				EnumValues: []string{"red", "green"},
			},
		},
	}

	result := compile.BuildEnums(spec, map[string]bool{"Color": true})

	assert.Empty(t, result.Enums)
}

func TestBuildEnums_NoEnums(t *testing.T) {
	spec := &compile.ParsedSpec{
		Schemas: []compile.ParsedSchema{
			{
				Name: "Widget",
				Properties: []compile.ParsedProperty{
					{Name: "name", Types: []string{"string"}},
				},
			},
		},
	}

	result := compile.BuildEnums(spec, nil)

	assert.Empty(t, result.Enums)
	assert.Empty(t, result.FieldRef)
}

func TestBuildEnums_MultipleInlineEnums(t *testing.T) {
	spec := &compile.ParsedSpec{
		Schemas: []compile.ParsedSchema{
			{
				Name: "Order",
				Properties: []compile.ParsedProperty{
					{Name: "status", Types: []string{"string"}, EnumValues: []string{"pending", "shipped"}},
					{Name: "priority", Types: []string{"string"}, EnumValues: []string{"low", "high"}},
				},
			},
		},
	}

	result := compile.BuildEnums(spec, nil)

	assert.Len(t, result.Enums, 2)
	assert.Contains(t, result.Enums, "OrderStatus")
	assert.Contains(t, result.Enums, "OrderPriority")
}
