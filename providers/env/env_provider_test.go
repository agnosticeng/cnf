package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	os.Clearenv()

	p := NewEnvProvider("APP")

	m, err := p.ReadMap()

	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{}, m)
}

func TestSimple(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP__CLIENTS__MY_SERVICE__READ_TIMEOUT", "100ms")

	p := NewEnvProvider("APP")

	m, err := p.ReadMap()

	assert.NoError(t, err)
	assert.Equal(
		t,
		map[string]interface{}{
			"Clients": map[string]interface{}{
				"MyService": map[string]interface{}{
					"ReadTimeout": "100ms",
				},
			},
		},
		m,
	)

	os.Clearenv()
}

func TestSlice(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP__CLUSTER__NAME", "test")
	os.Setenv("APP__CLUSTER__SERVERS__0__IP", "127.0.0.1")
	os.Setenv("APP__CLUSTER__SERVERS__1__IP", "127.0.0.2")
	os.Setenv("APP__CLUSTER__SERVERS__1__IFACES__0", "eth0")
	os.Setenv("APP__CLUSTER__SERVERS__1__IFACES__1", "eth1")

	p := NewEnvProvider("APP")

	m, err := p.ReadMap()

	assert.NoError(t, err)

	assert.Equal(
		t,
		map[string]interface{}{
			"Cluster": map[string]interface{}{
				"Name": "test",
				"Servers": []interface{}{
					map[string]interface{}{
						"Ip": "127.0.0.1",
					},
					map[string]interface{}{
						"Ip":     "127.0.0.2",
						"Ifaces": []interface{}{"eth0", "eth1"},
					},
				},
			},
		},
		m,
	)

	os.Clearenv()
}
