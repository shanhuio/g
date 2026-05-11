package dock

import (
	"testing"
)

func TestImageTag(t *testing.T) {
	for _, test := range []struct {
		s, img, tag string
	}{
		{s: "nextcloud", img: "nextcloud", tag: ""},
		{s: "nextcloud:19", img: "nextcloud", tag: "19"},
		{s: "shanhu.io/doorway", img: "shanhu.io/doorway", tag: ""},
		{s: "shanhu.io/doorway:v1", img: "shanhu.io/doorway", tag: "v1"},
	} {
		img, tag := ParseImageTag(test.s)
		if img != test.img || tag != test.tag {
			t.Errorf(
				"ParseImageTag(%q), got (%q, %q), want (%q, %q)",
				test.s, img, tag, test.img, test.tag,
			)
		}
	}
}
