package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	morpheyaml "github.com/kalo-build/morphe-go/pkg/yaml"
	"gopkg.in/yaml.v3"
)

func WriteStructures(outputPath string, structures map[string]morpheyaml.Structure) error {
	if len(structures) == 0 {
		return nil
	}
	structuresDir := filepath.Join(outputPath, "structures")
	if err := os.MkdirAll(structuresDir, 0755); err != nil {
		return fmt.Errorf("create structures dir: %w", err)
	}

	for _, name := range sortedKeys(structures) {
		structure := structures[name]
		data, err := marshalStructure(structure)
		if err != nil {
			return fmt.Errorf("marshal structure %s: %w", name, err)
		}
		fileName := SchemaNameToFileName(name) + ".str"
		filePath := filepath.Join(structuresDir, fileName)
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return fmt.Errorf("write structure %s: %w", name, err)
		}
	}
	return nil
}

func WriteEnums(outputPath string, enums map[string]morpheyaml.Enum) error {
	if len(enums) == 0 {
		return nil
	}
	enumsDir := filepath.Join(outputPath, "enums")
	if err := os.MkdirAll(enumsDir, 0755); err != nil {
		return fmt.Errorf("create enums dir: %w", err)
	}

	for _, name := range sortedKeys(enums) {
		enum := enums[name]
		data, err := marshalEnum(enum)
		if err != nil {
			return fmt.Errorf("marshal enum %s: %w", name, err)
		}
		fileName := SchemaNameToFileName(name) + ".enum"
		filePath := filepath.Join(enumsDir, fileName)
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return fmt.Errorf("write enum %s: %w", name, err)
		}
	}
	return nil
}

func marshalStructure(structure morpheyaml.Structure) ([]byte, error) {
	root := &yaml.Node{Kind: yaml.MappingNode}

	addScalar(root, "name", structure.Name)

	if len(structure.Fields) > 0 {
		fieldsNode := &yaml.Node{Kind: yaml.MappingNode}
		for _, name := range sortedKeys(structure.Fields) {
			f := structure.Fields[name]
			fieldNode := &yaml.Node{Kind: yaml.MappingNode}
			addScalar(fieldNode, "type", string(f.Type))
			if len(f.Attributes) > 0 {
				addStringSequence(fieldNode, "attributes", f.Attributes)
			}
			addMapping(fieldsNode, name, fieldNode)
		}
		addMapping(root, "fields", fieldsNode)
	}

	doc := &yaml.Node{
		Kind:    yaml.DocumentNode,
		Content: []*yaml.Node{root},
	}
	return yaml.Marshal(doc)
}

func marshalEnum(enum morpheyaml.Enum) ([]byte, error) {
	root := &yaml.Node{Kind: yaml.MappingNode}

	addScalar(root, "name", enum.Name)
	addScalar(root, "type", string(enum.Type))

	if len(enum.Entries) > 0 {
		entriesNode := &yaml.Node{Kind: yaml.MappingNode}
		keys := sortedKeys(enum.Entries)
		for _, k := range keys {
			addScalar(entriesNode, k, fmt.Sprintf("%v", enum.Entries[k]))
		}
		addMapping(root, "entries", entriesNode)
	}

	doc := &yaml.Node{
		Kind:    yaml.DocumentNode,
		Content: []*yaml.Node{root},
	}
	return yaml.Marshal(doc)
}

func addScalar(parent *yaml.Node, key string, value string) {
	parent.Content = append(parent.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: key},
		&yaml.Node{Kind: yaml.ScalarNode, Value: value},
	)
}

func addMapping(parent *yaml.Node, key string, value *yaml.Node) {
	parent.Content = append(parent.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: key},
		value,
	)
}

func addStringSequence(parent *yaml.Node, key string, values []string) {
	seqNode := &yaml.Node{Kind: yaml.SequenceNode, Style: yaml.FlowStyle}
	for _, v := range values {
		seqNode.Content = append(seqNode.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: v})
	}
	parent.Content = append(parent.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: key},
		seqNode,
	)
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
