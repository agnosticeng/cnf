package cnf

type Provider interface {
	ReadMap() (map[string]interface{}, error)
}
