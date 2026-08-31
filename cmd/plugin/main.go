package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kalo-build/plugin-openapi-morphe/pkg/compile"
	"github.com/kalo-build/plugin-openapi-morphe/pkg/compile/cfg"
)

type StoreConfig struct {
	ID        uint32 `json:"id"`
	Type      string `json:"type"`
	MountPath string `json:"mountPath,omitempty"`
}

type PluginConfig struct {
	Stores     map[string]StoreConfig `json:"stores,omitempty"`
	InputPath  string                 `json:"inputPath,omitempty"`
	OutputPath string                 `json:"outputPath,omitempty"`
	Config     PluginConfigFields     `json:"config"`
	Verbose    bool                   `json:"verbose,omitempty"`
}

type PluginConfigFields struct {
	SpecFileName  string   `json:"specFileName,omitempty"`
	SchemasToSkip []string `json:"schemasToSkip,omitempty"`
}

const (
	ErrMissingConfig      = 3
	ErrInvalidConfig      = 4
	ErrInputPathRequired  = 12
	ErrOutputPathRequired = 13
	ErrCompileFailed      = 1
)

func logInfo(verbose bool, format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stdout, format+"\n", args...)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: plugin-openapi-morphe <config>")
		os.Exit(ErrMissingConfig)
	}

	rawConfig := os.Args[1]
	var pluginConfig PluginConfig
	if err := json.Unmarshal([]byte(rawConfig), &pluginConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config JSON:", err)
		os.Exit(ErrInvalidConfig)
	}

	var inputPath, outputPath string

	if pluginConfig.Stores != nil {
		for _, store := range pluginConfig.Stores {
			switch store.MountPath {
			case "/input":
				inputPath = "/input"
			case "/output":
				outputPath = "/output"
			}
		}
	}

	if inputPath == "" && pluginConfig.InputPath != "" {
		inputPath = pluginConfig.InputPath
	}
	if outputPath == "" && pluginConfig.OutputPath != "" {
		outputPath = pluginConfig.OutputPath
	}

	if inputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Input path is required (directory containing OpenAPI spec)")
		os.Exit(ErrInputPathRequired)
	}
	if outputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Output path is required")
		os.Exit(ErrOutputPathRequired)
	}

	if inputPath[0] != '/' {
		if abs, err := filepath.Abs(inputPath); err == nil {
			inputPath = abs
		}
	}
	if outputPath[0] != '/' {
		if abs, err := filepath.Abs(outputPath); err == nil {
			outputPath = abs
		}
	}

	skipSet := make(map[string]bool, len(pluginConfig.Config.SchemasToSkip))
	for _, s := range pluginConfig.Config.SchemasToSkip {
		skipSet[s] = true
	}

	compileConfig := cfg.OpenAPIToMorpheConfig{
		InputDir:      inputPath,
		OutputPath:    outputPath,
		SpecFileName:  pluginConfig.Config.SpecFileName,
		SchemasToSkip: skipSet,
	}

	logInfo(pluginConfig.Verbose, "Reading OpenAPI spec from: '%s'", filepath.Join(compileConfig.InputDir, compileConfig.SpecFileName))
	logInfo(pluginConfig.Verbose, "Generating Morphe structures to: '%s'", compileConfig.OutputPath)

	if err := compile.OpenAPIToMorphe(compileConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Compilation failed:", err)
		os.Exit(ErrCompileFailed)
	}

	logInfo(pluginConfig.Verbose, "Morphe structures and enums generated successfully")
	os.Exit(0)
}
