# plugin-openapi-morphe

Generates Morphe structures (`.str`) and enums (`.enum`) from an OpenAPI 3.1 specification. Structures only — no models or entities (those require relationship inference that OpenAPI can't provide).

## What it generates

| OpenAPI artifact | Morphe output |
|------------------|---------------|
| **Component schema** (object with properties) | `.str` file — one structure per schema |
| **Inline string enum** (property with `enum`) | `.enum` file — extracted and named `{Schema}{Property}` |
| **Top-level string enum** (schema with only `enum`) | `.enum` file — keeps schema name |

### Example output

Given an OpenAPI schema:

```yaml
components:
  schemas:
    Customer:
      type: object
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
        email:
          type: string
        created_at:
          type: string
          format: date-time
      required: [id, name, email, created_at]
```

**Structure** (`customer.str`):

```yaml
name: Customer
fields:
    CreatedAt:
        type: Time
    Email:
        type: String
    ID:
        type: UUID
    Name:
        type: String
```

For inline enums:

```yaml
# OpenAPI
Invoice:
  type: object
  properties:
    status:
      type: string
      enum: [draft, open, paid, void]
```

**Enum** (`invoice_status.enum`):

```yaml
name: InvoiceStatus
type: String
entries:
    draft: draft
    open: open
    paid: paid
    void: void
```

### Type mappings

| OpenAPI type + format | Morphe type |
|-----------------------|-------------|
| `string` | `String` |
| `string` + `uuid` | `UUID` |
| `string` + `date-time` | `Time` |
| `string` + `date` | `Date` |
| `string` + `password` | `Protected` |
| `integer` | `Integer` |
| `number` | `Float` |
| `boolean` | `Boolean` |

### Property handling

| OpenAPI pattern | Morphe result |
|-----------------|---------------|
| `$ref: '#/components/schemas/X'` | Field type = `X` (structure reference) |
| `type: array` + `items.$ref` | Field type = referenced schema name |
| `enum: [...]` | Extracted as separate `.enum`, field type = enum name |
| Not in `required` | Field gets `optional` attribute |

## Input / output

| Direction | Format | Store suggestion | Description |
|-----------|--------|------------------|-------------|
| Input | `KA:OA1:YAML1` | `KA_OA_YAML` | OpenAPI 3.1 YAML specification |
| Output | `KA:MO1:YAML1` | `KA_MO_YAML` | Morphe YAML files (structures and enums) |

## Configuration

| Key | Type | Required | Default | Description |
|-----|------|----------|---------|-------------|
| `config.specFileName` | string | no | `"openapi.yaml"` | Name of the OpenAPI spec file to read |
| `config.schemasToSkip` | array | no | `[]` | Component schema names to exclude from generation |

## Pipeline context

```yaml
stores:
  KA_OA_YAML:
    format: "KA:OA1:YAML1"
    type: "localFileSystem"
    options:
      path: "./docs/openapi"

  KA_MO_YAML:
    format: "KA:MO1:YAML1"
    type: "localFileSystem"
    options:
      path: "./morphe"

plugins:
  "@kalo-build/plugin-openapi-morphe":
    version: "v1.0.0"
    inputs:
      openapi:
        format: "KA:OA1:YAML1"
        store: "KA_OA_YAML"
    output:
      format: "KA:MO1:YAML1"
      store: "KA_MO_YAML"
    config:
      specFileName: "openapi.yaml"
      schemasToSkip: []
```

## Project structure

```
plugin-openapi-morphe/
├── cmd/plugin/             # WASM entry point
├── pkg/
│   └── compile/            # Compilation pipeline
│       ├── compile.go      # OpenAPIToMorphe entry point
│       ├── parse.go        # OpenAPI spec parsing (libopenapi)
│       ├── build_structures.go  # Schema → Morphe structure
│       ├── build_enums.go  # Inline enum extraction
│       ├── typemap.go      # OpenAPI → Morphe type mapping
│       ├── naming.go       # Name conversion utilities
│       ├── write.go        # Morphe YAML file writers
│       └── cfg/            # Configuration types
├── internal/testutils/     # Test path helpers
├── testdata/
│   ├── input/              # Sample OpenAPI spec
│   └── ground-truth/       # Expected Morphe output
├── dist/                   # WASM output
├── scripts/                # Build scripts
└── plugin.yaml             # Kalo plugin manifest
```

## Building

```bash
# Native binary
go build ./cmd/plugin

# WASM (for Kalo CLI)
GOOS=wasip1 GOARCH=wasm go build -o dist/plugin.wasm cmd/plugin/main.go
```

## Testing

```bash
go test ./...
```
