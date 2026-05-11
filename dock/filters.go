package dock

import (
	"encoding/json"

	"shanhu.io/g/errcode"
)

func labelFilters(label string) (string, error) {
	filters := map[string][]string{"label": {label}}
	filterBytes, err := json.Marshal(filters)
	if err != nil {
		return "", errcode.Annotate(err, "marshal filter")
	}
	return string(filterBytes), nil
}
