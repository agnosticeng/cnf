package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type FileProvider struct {
	path string
}

func NewFileProvider(path string) *FileProvider {
	return &FileProvider{
		path: path,
	}
}

func (p *FileProvider) ReadMap() (map[string]interface{}, error) {
	content, err := os.ReadFile(p.path)

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var unmarshaler func([]byte, interface{}) error

	switch filepath.Ext(p.path) {
	case ".json":
		unmarshaler = json.Unmarshal
	case ".yaml", ".yml":
		unmarshaler = yaml.Unmarshal
	default:
		return nil, fmt.Errorf("unhandled file extension: %s", filepath.Ext(p.path))
	}

	var m map[string]interface{}

	if err := unmarshaler(content, &m); err != nil {
		return nil, err
	}

	return m, nil
}
