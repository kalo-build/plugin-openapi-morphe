package compile_test

import (
	"testing"

	"github.com/kalo-build/plugin-openapi-morphe/pkg/compile"
	"github.com/stretchr/testify/assert"
)

func TestPropertyNameToFieldName_SnakeCase(t *testing.T) {
	assert.Equal(t, "CreatedAt", compile.PropertyNameToFieldName("created_at"))
}

func TestPropertyNameToFieldName_SnakeCaseWithID(t *testing.T) {
	assert.Equal(t, "CustomerID", compile.PropertyNameToFieldName("customer_id"))
}

func TestPropertyNameToFieldName_CamelCase(t *testing.T) {
	assert.Equal(t, "CreatedAt", compile.PropertyNameToFieldName("createdAt"))
}

func TestPropertyNameToFieldName_AlreadyPascal(t *testing.T) {
	assert.Equal(t, "Name", compile.PropertyNameToFieldName("Name"))
}

func TestPropertyNameToFieldName_SingleWord(t *testing.T) {
	assert.Equal(t, "Email", compile.PropertyNameToFieldName("email"))
}

func TestPropertyNameToFieldName_Abbreviations(t *testing.T) {
	assert.Equal(t, "APIURL", compile.PropertyNameToFieldName("api_url"))
}

func TestPropertyNameToFieldName_ID(t *testing.T) {
	assert.Equal(t, "ID", compile.PropertyNameToFieldName("id"))
}

func TestSchemaNameToFileName_PascalCase(t *testing.T) {
	assert.Equal(t, "customer_create", compile.SchemaNameToFileName("CustomerCreate"))
}

func TestSchemaNameToFileName_SingleWord(t *testing.T) {
	assert.Equal(t, "customer", compile.SchemaNameToFileName("Customer"))
}

func TestSchemaNameToFileName_WithAcronym(t *testing.T) {
	assert.Equal(t, "api_key", compile.SchemaNameToFileName("APIKey"))
}

func TestRefToSchemaName(t *testing.T) {
	assert.Equal(t, "Customer", compile.RefToSchemaName("#/components/schemas/Customer"))
}

func TestRefToSchemaName_Nested(t *testing.T) {
	assert.Equal(t, "InvoiceItem", compile.RefToSchemaName("#/components/schemas/InvoiceItem"))
}
