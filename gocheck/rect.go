package gocheck

import (
	"os"

	"shanhu.io/g/lexing"
	"shanhu.io/g/textbox"
)

// CheckRect checks if all the files are within the given rectangle.
func CheckRect(files []string, h, w int) []*lexing.Error {
	errs := lexing.NewErrorList()
	for _, f := range files {
		fin, err := os.Open(f)
		if lexing.LogError(errs, err) {
			continue
		}

		errs.AddAll(textbox.CheckRect(f, fin, h, w))
		if lexing.LogError(errs, fin.Close()) {
			continue
		}
	}

	return errs.Errs()
}
