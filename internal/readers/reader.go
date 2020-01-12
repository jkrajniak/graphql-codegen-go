package readers

type SchemaReader interface {
	Read() ([]byte, error)
}


func DiscoverReader(schemaPath string) SchemaReader {
	return NewLocalReader(schemaPath)
}