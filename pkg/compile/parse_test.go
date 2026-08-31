package compile_test

import (
	"testing"

	"github.com/kalo-build/plugin-openapi-morphe/pkg/compile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const minimalSpec = `
openapi: "3.1.0"
info:
  title: Test
  version: 1.0.0
paths: {}
components:
  schemas:
    Widget:
      type: object
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
        count:
          type: integer
      required:
        - id
        - name
`

func TestParseOpenAPISpec_MinimalSchema(t *testing.T) {
	spec, err := compile.ParseOpenAPISpec([]byte(minimalSpec))
	require.NoError(t, err)
	require.Len(t, spec.Schemas, 1)

	schema := spec.Schemas[0]
	assert.Equal(t, "Widget", schema.Name)
	assert.Len(t, schema.Properties, 3)
	assert.True(t, schema.Required["id"])
	assert.True(t, schema.Required["name"])
	assert.False(t, schema.Required["count"])
}

func TestParseOpenAPISpec_PropertyTypes(t *testing.T) {
	spec, err := compile.ParseOpenAPISpec([]byte(minimalSpec))
	require.NoError(t, err)

	schema := spec.Schemas[0]
	propByName := make(map[string]compile.ParsedProperty)
	for _, p := range schema.Properties {
		propByName[p.Name] = p
	}

	idProp := propByName["id"]
	assert.Equal(t, []string{"string"}, idProp.Types)
	assert.Equal(t, "uuid", idProp.Format)

	nameProp := propByName["name"]
	assert.Equal(t, []string{"string"}, nameProp.Types)
	assert.Empty(t, nameProp.Format)

	countProp := propByName["count"]
	assert.Equal(t, []string{"integer"}, countProp.Types)
}

const refSpec = `
openapi: "3.1.0"
info:
  title: Test
  version: 1.0.0
paths: {}
components:
  schemas:
    Order:
      type: object
      properties:
        customer:
          $ref: '#/components/schemas/Customer'
      required:
        - customer
    Customer:
      type: object
      properties:
        name:
          type: string
      required:
        - name
`

func TestParseOpenAPISpec_RefProperty(t *testing.T) {
	spec, err := compile.ParseOpenAPISpec([]byte(refSpec))
	require.NoError(t, err)
	require.Len(t, spec.Schemas, 2)

	orderSchema := spec.Schemas[1]
	assert.Equal(t, "Order", orderSchema.Name)

	customerProp := orderSchema.Properties[0]
	assert.Equal(t, "customer", customerProp.Name)
	assert.Equal(t, "#/components/schemas/Customer", customerProp.Ref)
}

const enumSpec = `
openapi: "3.1.0"
info:
  title: Test
  version: 1.0.0
paths: {}
components:
  schemas:
    Task:
      type: object
      properties:
        status:
          type: string
          enum:
            - open
            - closed
      required:
        - status
`

func TestParseOpenAPISpec_InlineEnum(t *testing.T) {
	spec, err := compile.ParseOpenAPISpec([]byte(enumSpec))
	require.NoError(t, err)

	taskSchema := spec.Schemas[0]
	statusProp := taskSchema.Properties[0]
	assert.Equal(t, "status", statusProp.Name)
	assert.Equal(t, []string{"open", "closed"}, statusProp.EnumValues)
}

const topLevelEnumSpec = `
openapi: "3.1.0"
info:
  title: Test
  version: 1.0.0
paths: {}
components:
  schemas:
    Color:
      type: string
      enum:
        - red
        - green
        - blue
`

func TestParseOpenAPISpec_TopLevelEnum(t *testing.T) {
	spec, err := compile.ParseOpenAPISpec([]byte(topLevelEnumSpec))
	require.NoError(t, err)
	require.Len(t, spec.Schemas, 1)

	schema := spec.Schemas[0]
	assert.Equal(t, "Color", schema.Name)
	assert.Empty(t, schema.Properties)
	assert.Equal(t, []string{"red", "green", "blue"}, schema.EnumValues)
}

func TestParseOpenAPISpec_InvalidSpec(t *testing.T) {
	_, err := compile.ParseOpenAPISpec([]byte("not valid yaml: [[["))
	assert.Error(t, err)
}

const arrayRefSpec = `
openapi: "3.1.0"
info:
  title: Test
  version: 1.0.0
paths: {}
components:
  schemas:
    OrderList:
      type: object
      properties:
        items:
          type: array
          items:
            $ref: '#/components/schemas/Order'
      required:
        - items
    Order:
      type: object
      properties:
        id:
          type: string
      required:
        - id
`

func TestParseOpenAPISpec_ArrayWithRef(t *testing.T) {
	spec, err := compile.ParseOpenAPISpec([]byte(arrayRefSpec))
	require.NoError(t, err)

	orderListSchema := spec.Schemas[1]
	assert.Equal(t, "OrderList", orderListSchema.Name)

	itemsProp := orderListSchema.Properties[0]
	assert.Equal(t, "items", itemsProp.Name)
	assert.True(t, itemsProp.IsArray)
	assert.Equal(t, "#/components/schemas/Order", itemsProp.ItemRef)
}
