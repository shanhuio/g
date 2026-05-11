package dock

import (
	"path"
	"strings"
)

// ParseImageTag parses "image:tag" into ("image", "tag").
func ParseImageTag(s string) (image, tag string) {
	base := path.Base(s)
	colon := strings.LastIndex(base, ":")
	if colon < 0 {
		return s, ""
	}

	t := base[colon+1:]
	img := strings.TrimSuffix(s, base[colon:])
	return img, t
}
