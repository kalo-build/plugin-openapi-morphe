package compile_test

import (
	"testing"

	morpheyaml "github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-openapi-morphe/pkg/compile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildStructures_BasicFields(t *testing.T) {
	spec := &compile.ParsedSpec{
		Schemas: []compile.ParsedSchema{
			{
				Name: "Widget",
				Properties: []compile.ParsedProperty{
					{Name: "id", Types: []string{"string"}, Format: "uuid", Required: true},
					{Name: "name", Types: []string{"string"}, Required: true},
					{Name: "count", Types: []string{"integer"}, Required: false},
				},
				Required: map[string]bool{"id": true, "name": true},
			},
		},
	}
	enumRefs := &compile.BuildEnumsResult{
		Enums:    make(map[string]morpheyaml.Enum),
		FieldRef: make(map[string]string),
	}

	structures := compile.BuildStructures(spec, enumRefs, nil)

	require.Contains(t, structures, "Widget")
	widget := structures["Widget"]
	assert.Equal(t, "Widget", widget.Name)

	assert.Equal(t, morpheyaml.StructureFieldTypeUUID, widget.Fields["ID"].Type)
	assert.Empty(t, widget.Fields["ID"].Attributes)

	assert.Equal(t, morpheyaml.StructureFieldTypeString, widget.Fields["Name"].Type)
	assert.Empty(t, widget.Fields["Name"].Attributes)

	assert.Equal(t, morpheyaml.StructureFieldTypeInteger, widget.Fields["Count"].Type)
	assert.Equal(t, []string{"optional"}, widget.Fields["Count"].Attributes)
}

func TestBuildStructures_RefField(t *testing.T) {
	spec := &compile.ParsedSpec{
		Schemas: []compile.ParsedSchema{
			{
				Name: "Order",
				Properties: []compile.ParsedProperty{
					{Name: "customer", Ref: "#/components/schemas/Customer", Required: true},
				},
				Required: map[string]bool{"customer": true},
			},
		},
	}
	enumRefs := &compile.BuildEnumsResult{
		Enums:    make(map[string]morpheyaml.Enum),
		FieldRef: make(map[string]string),
	}

	structures := compile.BuildStructures(spec, enumRefs, nil)

	require.Contains(t, structures, "Order")
	assert.Equal(t, morpheyaml.StructureFieldType("Customer"), structures["Order"].Fields["Customer"].Type)
}

func TestBuildStructures_EnumRefField(t *testing.T) {
	spec := &compile.ParsedSpec{
		Schemas: []compile.ParsedSchema{
			{
				Name: "Task",
				Properties: []compile.ParsedProperty{
					{Name: "status", Types: []string{"string"}, EnumValues: []string{"open", "closed"}, Required: true},
				},
				Required: map[string]bool{"status": true},
			},
		},
	}
	enumRefs := &compile.BuildEnumsResult{
		Enums: map[string]morpheyaml.Enum{
			"TaskStatus": {Name: "TaskStatus", Type: morpheyaml.EnumTypeString},
		},
		FieldRef: map[string]string{
			"Task.status": "TaskStatus",
		},
	}

	structures := compile.BuildStructures(spec, enumRefs, nil)

	require.Contains(t, structures, "Task")
	assert.Equal(t, morpheyaml.StructureFieldType("TaskStatus"), structures["Task"].Fields["Status"].Type)
}

func TestBuildStructures_SkipsEnumOnlySchemas(t *testing.T) {
	spec := &compile.ParsedSpec{
		Schemas: []compile.ParsedSchema{
			{
				Name:       "Color",
				EnumValues: []string{"red", "blue"},
			},
		},
	}
	enumRefs := &compile.BuildEnumsResult{
		Enums:    make(map[string]morpheyaml.Enum),
		FieldRef: make(map[string]string),
	}

	structures := compile.BuildStructures(spec, enumRefs, nil)

	assert.Empty(t, structures)
}

func TestBuildStructures_SkippedSchema(t *testing.T) {
	spec := &compile.ParsedSpec{
		Schemas: []compile.ParsedSchema{
			{
				Name: "Internal",
				Properties: []compile.ParsedProperty{
					{Name: "id", Types: []string{"string"}, Required: true},
				},
			},
		},
	}
	enumRefs := &compile.BuildEnumsResult{
		Enums:    make(map[string]morpheyaml.Enum),
		FieldRef: make(map[string]string),
	}

	structures := compile.BuildStructures(spec, enumRefs, map[string]bool{"Internal": true})

	assert.Empty(t, structures)
}

func TestBuildStructures_ArrayRefField(t *testing.T) {
	spec := &compile.ParsedSpec{
		Schemas: []compile.ParsedSchema{
			{
				Name: "OrderList",
				Properties: []compile.ParsedProperty{
					{Name: "items", IsArray: true, ItemRef: "#/components/schemas/Order", Required: true},
				},
				Required: map[string]bool{"items": true},
			},
		},
	}
	enumRefs := &compile.BuildEnumsResult{
		Enums:    make(map[string]morpheyaml.Enum),
		FieldRef: make(map[string]string),
	}

	structures := compile.BuildStructures(spec, enumRefs, nil)

	require.Contains(t, structures, "OrderList")
	assert.Equal(t, morpheyaml.StructureFieldType("Order"), structures["OrderList"].Fields["Items"].Type)
}

func TestBuildStructures_OptionalField(t *testing.T) {
	spec := &compile.ParsedSpec{
		Schemas: []compile.ParsedSchema{
			{
				Name: "Profile",
				Properties: []compile.ParsedProperty{
					{Name: "bio", Types: []string{"string"}, Required: false},
				},
			},
		},
	}
	enumRefs := &compile.BuildEnumsResult{
		Enums:    make(map[string]morpheyaml.Enum),
		FieldRef: make(map[string]string),
	}

	structures := compile.BuildStructures(spec, enumRefs, nil)

	require.Contains(t, structures, "Profile")
	bio := structures["Profile"].Fields["Bio"]
	assert.Equal(t, []string{"optional"}, bio.Attributes)
}
