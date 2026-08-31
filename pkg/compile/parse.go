package compile

import (
	"fmt"
	"sort"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type ParsedSpec struct {
	Schemas []ParsedSchema
}

type ParsedSchema struct {
	Name       string
	Properties []ParsedProperty
	Required   map[string]bool
	EnumValues []string
}

type ParsedProperty struct {
	Name       string
	Types      []string
	Format     string
	Ref        string
	Required   bool
	EnumValues []string
	IsArray    bool
	ItemRef    string
	ItemTypes  []string
	ItemFormat string
}

func ParseOpenAPISpec(specBytes []byte) (*ParsedSpec, error) {
	document, err := libopenapi.NewDocument(specBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI document: %w", err)
	}

	model, errs := document.BuildV3Model()
	if errs != nil {
		return nil, fmt.Errorf("failed to build V3 model: %w", errs)
	}

	return extractSchemas(&model.Model), nil
}

func extractSchemas(v3doc *v3.Document) *ParsedSpec {
	spec := &ParsedSpec{}

	if v3doc.Components == nil || v3doc.Components.Schemas == nil {
		return spec
	}

	for pair := v3doc.Components.Schemas.First(); pair != nil; pair = pair.Next() {
		name := pair.Key()
		proxy := pair.Value()
		schema := convertToParsedSchema(name, proxy)
		spec.Schemas = append(spec.Schemas, schema)
	}

	sort.Slice(spec.Schemas, func(i, j int) bool {
		return spec.Schemas[i].Name < spec.Schemas[j].Name
	})

	return spec
}

func convertToParsedSchema(name string, proxy *base.SchemaProxy) ParsedSchema {
	schema := ParsedSchema{
		Name:     name,
		Required: make(map[string]bool),
	}

	if proxy == nil {
		return schema
	}
	resolved := proxy.Schema()
	if resolved == nil {
		return schema
	}

	for _, r := range resolved.Required {
		schema.Required[r] = true
	}

	if resolved.Enum != nil {
		for _, e := range resolved.Enum {
			schema.EnumValues = append(schema.EnumValues, fmt.Sprintf("%v", e.Value))
		}
	}

	if resolved.Properties != nil {
		for pair := resolved.Properties.First(); pair != nil; pair = pair.Next() {
			propName := pair.Key()
			propProxy := pair.Value()
			prop := convertToParsedProperty(propName, propProxy, schema.Required[propName])
			schema.Properties = append(schema.Properties, prop)
		}
		sort.Slice(schema.Properties, func(i, j int) bool {
			return schema.Properties[i].Name < schema.Properties[j].Name
		})
	}

	return schema
}

func convertToParsedProperty(name string, proxy *base.SchemaProxy, required bool) ParsedProperty {
	prop := ParsedProperty{
		Name:     name,
		Required: required,
	}

	if proxy == nil {
		return prop
	}

	ref := proxy.GetReference()
	if ref != "" {
		prop.Ref = ref
		return prop
	}

	resolved := proxy.Schema()
	if resolved == nil {
		return prop
	}

	prop.Types = resolved.Type
	prop.Format = resolved.Format

	if resolved.Enum != nil {
		for _, e := range resolved.Enum {
			prop.EnumValues = append(prop.EnumValues, fmt.Sprintf("%v", e.Value))
		}
	}

	if hasType(resolved.Type, "array") && resolved.Items != nil && resolved.Items.A != nil {
		prop.IsArray = true
		itemRef := resolved.Items.A.GetReference()
		if itemRef != "" {
			prop.ItemRef = itemRef
		} else {
			itemSchema := resolved.Items.A.Schema()
			if itemSchema != nil {
				prop.ItemTypes = itemSchema.Type
				prop.ItemFormat = itemSchema.Format
			}
		}
	}

	return prop
}

func hasType(types []string, t string) bool {
	for _, tt := range types {
		if tt == t {
			return true
		}
	}
	return false
}
