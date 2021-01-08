package codec

import (
	"log"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/encoding"
)

func init() {
	log.Print("test")
	encoding.RegisterCodec(codec{})
}

type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	log.Print(v)
	return proto.Marshal(v.(proto.Message))
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	log.Print(data)
	return proto.Unmarshal(data, v.(proto.Message))
}

func (codec) Name() string {
	return "proto-plus"
}
