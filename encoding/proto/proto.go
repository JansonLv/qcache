// Package proto defines the protobuf codec. Importing this package will
// register the codec.
package proto

import (
	"errors"
	"github.com/JansonLv/qcache/encoding"
	"google.golang.org/protobuf/proto"
)

// Name is the name registered for the proto compressor.
const Name = "proto"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with protobuf. It is the default codec for Transport.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	if _, ok := v.(proto.Message); !ok {
		return nil, errors.New("value is not a proto.Message")
	}
	return proto.Marshal(v.(proto.Message))
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	if _, ok := v.(proto.Message); !ok {
		return errors.New("value is not a proto.Message")
	}
	return proto.Unmarshal(data, v.(proto.Message))
}

func (codec) Copy(dst interface{}, src interface{}) error {
	to, ok := dst.(proto.Message)
	if !ok {
		return errors.New("value is not a proto.Message")
	}
	form, ok := src.(proto.Message)
	if !ok {
		return errors.New("value is not a proto.Message")
	}
	proto.Merge(to, form)
	return nil
}

func (codec) Name() string {
	return Name
}
