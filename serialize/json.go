package serialize

import "encoding/json"

type jsonSerializer struct{}

func (s jsonSerializer) Serialize(v any) ([]byte, error) {
	return json.Marshal(v)
}
