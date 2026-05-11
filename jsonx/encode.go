package jsonx

import (
	"encoding/json"
	"fmt"
	"io"

	"shanhu.io/g/errcode"
)

func writeString(w io.Writer, s string) error {
	_, err := io.WriteString(w, s)
	return err
}

func encodeJSON(w io.Writer, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	return err
}

func encodeBasic(w io.Writer, v *basic) error {
	if v.lead != nil {
		if v.lead.Type == tokOperator && v.lead.Lit == "-" {
			if err := writeString(w, "-"); err != nil {
				return err
			}
		}
	}

	switch v.token.Type {
	case tokInt:
		if err := writeString(w, v.token.Lit); err != nil {
			return err
		}
	case tokFloat, tokString:
		if err := encodeJSON(w, v.value); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected token: %s", tokenTypeStr(v.token))
	}
	return nil
}

func encodeObject(w io.Writer, v *object) error {
	if err := writeString(w, "{"); err != nil {
		return err
	}
	for i, entry := range v.entries {
		if i > 0 {
			if err := writeString(w, ","); err != nil {
				return err
			}
		}
		k := entry.key
		if k.token.Type == tokIdent {
			if err := encodeJSON(w, k.token.Lit); err != nil {
				return err
			}
		} else {
			if err := encodeJSON(w, k.value); err != nil {
				return err
			}
		}
		if err := writeString(w, ":"); err != nil {
			return err
		}
		if err := encodeValue(w, entry.value); err != nil {
			return err
		}
	}
	return writeString(w, "}")
}

func encodeList(w io.Writer, v *list) error {
	if err := writeString(w, "["); err != nil {
		return err
	}
	for i, entry := range v.entries {
		if i > 0 {
			if err := writeString(w, ","); err != nil {
				return err
			}
		}
		if err := encodeValue(w, entry.value); err != nil {
			return err
		}
	}
	return writeString(w, "]")
}

func encodeIdentList(w io.Writer, v *identList) error {
	if err := writeString(w, "["); err != nil {
		return err
	}
	for i, entry := range v.entries {
		if i > 0 {
			if err := writeString(w, ","); err != nil {
				return err
			}
		}
		if err := encodeJSON(w, entry.Lit); err != nil {
			return err
		}
	}
	return writeString(w, "]")
}

func encodeValue(w io.Writer, v value) error {
	switch v := v.(type) {
	case *null:
		if err := writeString(w, "null"); err != nil {
			return err
		}
	case *basic:
		if err := encodeBasic(w, v); err != nil {
			return err
		}
	case *boolean:
		if err := writeString(w, v.keyword.Lit); err != nil {
			return err
		}
	case *object:
		if err := encodeObject(w, v); err != nil {
			return err
		}
	case *list:
		if err := encodeList(w, v); err != nil {
			return err
		}
	case *identList:
		if err := encodeIdentList(w, v); err != nil {
			return err
		}
	default:
		return errcode.Internalf("invalid type: %T", v)
	}
	return nil
}
