package cnf

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/agnosticeng/cnf/providers/env"
	"github.com/agnosticeng/cnf/providers/file"
	"github.com/stretchr/testify/assert"
)

type config struct {
	Amount  int
	Tag     string
	Name    string
	Timeout time.Duration
	Sub     subConfig
}

type subConfig struct {
	Impl   string
	Target *url.URL
	Items  []subConfigItem
}

type subConfigItem struct {
	Host              string
	Port              int
	ConnectionTimeout time.Duration
	RequestTimeout    time.Duration
}

func defaultConfig() config {
	return config{
		Amount:  1,
		Tag:     "lolo",
		Timeout: time.Second,
		Sub: subConfig{
			Impl: "coco",
			Items: []subConfigItem{
				{
					Host:              "test.local",
					Port:              80,
					ConnectionTimeout: time.Second,
					RequestTimeout:    time.Second,
				},
			},
		},
	}
}

func TestLoad(t *testing.T) {
	os.Setenv("APP__NAME", "TEST")
	os.Setenv("APP__SUB__TARGET", "s3://my-bucket/coco.json")
	os.Setenv("APP__SUB__ITEMS__0__PORT", "443")
	os.Setenv("APP__SUB__ITEMS__1__HOST", "loco.local")
	os.Setenv("APP__SUB__ITEMS__1__PORT", "8080")
	os.Setenv("APP__SUB__ITEMS__1__REQUEST_TIMEOUT", "17s")

	var cfg = defaultConfig()

	if err := Load(
		&cfg,
		WithProvider(file.NewFileProvider("examples/basic/conf.yaml")),
		WithProvider(env.NewEnvProvider("APP")),
	); err != nil {
		panic(err)
	}

	assert.Equal(
		t,
		config{
			Amount:  1,
			Tag:     "sample",
			Name:    "TEST",
			Timeout: time.Second * 5,
			Sub: subConfig{
				Impl: "coco",
				Target: &url.URL{
					Scheme: "s3",
					Host:   "my-bucket",
					Path:   "/coco.json",
				},
				Items: []subConfigItem{
					{
						Host:              "test.local",
						Port:              443,
						ConnectionTimeout: time.Second,
						RequestTimeout:    time.Second,
					},
					{
						Host:           "loco.local",
						Port:           8080,
						RequestTimeout: time.Second * 17,
					},
				},
			},
		},
		cfg,
	)
}
