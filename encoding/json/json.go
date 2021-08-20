package json

import (
	"encoding/json"
	"github.com/JansonLv/qcache/encoding"
	"github.com/jinzhu/copier"
)

// Name is the name registered for the json codec.
const Name = "json"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with json.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (codec) Copy(dst interface{}, src interface{}) error {
	return copier.Copy(dst, src)
}

func (codec) Name() string {
	return Name
}
