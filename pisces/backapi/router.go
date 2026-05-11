package backapi

import (
	"shanhu.io/g/aries"
	"shanhu.io/g/pisces"
)

// Router provides an API service router for the given PsqlTables.
func Router(b *pisces.Tables) *aries.Router {
	r := aries.NewRouter()

	r.File("create", func(*aries.C) error { return b.Create() })
	r.File("create-missing", func(*aries.C) error {
		return b.CreateMissing()
	})
	r.File("destroy", func(*aries.C) error { return b.Destroy() })

	return r
}
