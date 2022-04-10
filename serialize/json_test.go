package serialize_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/desotech-it/whoami/serialize"
)

func TestNilSerializesToNull(t *testing.T) {
	data, err := serialize.SerializerJSON.Serialize(nil)
	if err != nil {
		t.Error(err)
	}
	got := data
	want := []byte{'n', 'u', 'l', 'l'}
	if !bytes.Equal(got, want) {
		t.Errorf("got = %v; want = %v", got, want)
	}
}

func TestNonNilSerializesToPayload(t *testing.T) {
	type payload struct {
		Str string
		Num int
	}
	p := payload{"foo", 42}
	data, err := serialize.SerializerJSON.Serialize(p)
	if err != nil {
		t.Error(err)
	}
	var got payload
	want := p
	if err := json.Unmarshal(data, &got); err != nil {
		t.Error(err)
	}
	if want != got {
		t.Errorf("got = %v; want = %v", got, want)
	}
}
