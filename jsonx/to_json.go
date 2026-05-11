package jsonx

import (
	"bytes"

	"shanhu.io/g/lexing"
)

func marshalValue(v value) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := encodeValue(buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func marshalValueLexing(v value) ([]byte, []*lexing.Error) {
	bs, err := marshalValue(v)
	if err != nil {
		return nil, lexing.SingleErr(err)
	}
	return bs, nil
}

// ToJSON converts a JSONX stream into a JSON stream.
func ToJSON(input []byte) ([]byte, []*lexing.Error) {
	r := bytes.NewReader(input)
	p, _ := newParser("", r)
	v := parseValue(p)
	if errs := p.Errs(); errs != nil {
		return nil, errs
	}
	return marshalValueLexing(v)
}
