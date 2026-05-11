package pathutil

import (
	"testing"
)

func TestValidPath(t *testing.T) {
	o := func(p string) {
		if !ValidPath(p) {
			t.Errorf("%q shoud be valid", p)
		}
	}

	e := func(p string) {
		if ValidPath(p) {
			t.Errorf("%q should be invalid", p)
		}
	}

	o("/")
	o("/asdf")
	o("/valentines_day")
	o("/thank/you")
	o("/3307")
	o("/c323/b75/53_df_")

	e("")
	e("/Hello")
	e("//")
	e("/as/")
	e("/as//of")
	e("/asdf-er")
	e("/  ")
	e("asdf")
	e("/2014-01-18")
}
