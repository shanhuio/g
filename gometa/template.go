// Copyright (C) 2023  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package gometa

import (
	"html/template"
)

var tmpl = template.Must(template.New("index").Parse(`<!doctype html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
	<meta name="go-import" content="{{.ImportRoot}} {{.VCS}} {{.VCSRoot}}">
</head>
<body>{{.ImportRoot}}</body>
</html>
`))

type data struct {
	ImportRoot string
	VCS        string
	VCSRoot    string
	Suffix     string
}
