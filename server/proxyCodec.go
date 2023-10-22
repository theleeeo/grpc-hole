package server

import "google.golang.org/grpc/encoding"

func init() {
	encoding.RegisterCodec(proxyCodec{})
}

type proxyCodec struct{}

func (proxyCodec) Marshal(v interface{}) ([]byte, error) {
	return *(v.(*[]byte)), nil
}

func (proxyCodec) Unmarshal(data []byte, v interface{}) error {
	vv := v.(*[]byte)
	*vv = data
	return nil
}

func (proxyCodec) Name() string {
	return "proxy"
}
