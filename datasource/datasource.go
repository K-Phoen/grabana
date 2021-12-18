package datasource

type Datasource interface {
	Name() string
	MarshalJSON() ([]byte, error)
}
