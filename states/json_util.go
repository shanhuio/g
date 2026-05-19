package states

import (
	"encoding/json"

	"shanhu.io/g/jsonx"
)

// GetJSON gets a JSON encoded state.
func GetJSON(ctx C, s States, key string, v any) error {
	bs, err := s.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}

// GetJSONX gets a JSONX encoded state.
func GetJSONX(ctx C, s States, key string, v any) error {
	bs, err := s.Get(ctx, key)
	if err != nil {
		return err
	}
	return jsonx.Unmarshal(bs, v)
}

// PutJSON puts a JSON encoded state.
func PutJSON(ctx C, s States, key string, v any) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.Put(ctx, key, bs)
}
