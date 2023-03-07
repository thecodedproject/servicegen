package resources

//go:generate resourcegen --struct_name=resources
type resources struct {
}

func New() *resources {

	return &resources{}
}
