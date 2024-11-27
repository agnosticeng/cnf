package env

import (
	"os"
	"sort"
	"strings"

	"github.com/agnosticeng/dynamap"
	"github.com/iancoleman/strcase"
)

type EnvProvider struct {
	prefix string
}

func NewEnvProvider(prefix string) *EnvProvider {
	return &EnvProvider{
		prefix: prefix,
	}
}

func (p *EnvProvider) ReadMap() (map[string]interface{}, error) {
	var (
		res interface{} = make(map[string]interface{})
		err error
		env = make([]string, len(os.Environ()))
	)

	copy(env, os.Environ())
	sort.Strings(env)

	for _, kv := range env {
		segments := strings.SplitN(kv, "=", 2)
		path := strings.Split(segments[0], "__")

		if len(p.prefix) > 0 {
			if path[0] != p.prefix {
				continue
			} else {
				path = path[1:]
			}
		}

		for i := 0; i < len(path); i++ {
			path[i] = strings.ToLower(path[i])
			path[i] = strcase.ToCamel(path[i])
		}

		res, err = dynamap.SSet(res, segments[1], path...)

		if err != nil {
			return nil, err
		}
	}

	return res.(map[string]interface{}), nil
}
