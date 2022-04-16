package serialize

type Serializer interface {
	Serialize(v any) ([]byte, error)
}

var (
	SerializerJSON Serializer = jsonSerializer{}
)
