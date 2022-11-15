package suprsend

import (
	"embed"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed request_json
var fs embed.FS

var schemaCache = map[string]*gojsonschema.Schema{}

/*
Returns schema from memory cache. If not already in memory, loads it from the file system.
Returns error if either schema-file is not present or has invalid jsonschema format
*/
func GetSchema(schemaName string) (*gojsonschema.Schema, error) {
	if _, found := schemaCache[schemaName]; !found {
		schema, err := loadJsonSchema(schemaName)
		if err != nil {
			return nil, err
		}
		schemaCache[schemaName] = schema
	}
	return schemaCache[schemaName], nil
}

func loadJsonSchema(schemaName string) (*gojsonschema.Schema, error) {
	// Get relative path of the jsonschema file
	relPath := fmt.Sprintf("request_json/%s.json", schemaName)
	content, err := fs.ReadFile(relPath)
	if err != nil {
		return nil, fmt.Errorf("SuprsendMissingSchema: %s. %w", schemaName, err)
	}
	if len(content) == 0 {
		return nil, fmt.Errorf("SuprsendMissingSchema: %s. %v", schemaName, "empty content")
	}
	// loading the schema
	var res gojsonschema.JSONLoader = gojsonschema.NewBytesLoader(content)
	schema, err := gojsonschema.NewSchema(res)
	if err != nil {
		return nil, fmt.Errorf("SuprsendInvalidSchema: %s. %w", schemaName, err)
	}
	return schema, nil
}
