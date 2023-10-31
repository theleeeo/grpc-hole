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
	vv := v.(*[]byte) //nolint:errcheck // If there even could be an error here it might as well just panic
	*vv = data
	return nil
}

func (proxyCodec) Name() string {
	return "proxy"
}
