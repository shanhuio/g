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

package webgen

// NewHTML creates a new blank HTML element with the specified language.
func NewHTML(lang string) *Node {
	if lang == "" {
		return HTML()
	}
	return HTML(Attrs{"lang": lang})
}

// NewHTMLEnglish creates a new English HTML element.
func NewHTMLEnglish() *Node { return NewHTML("en") }

// NewHTMLChinese creates a new Chinese HTML element.
func NewHTMLChinese() *Node { return NewHTML("zh") }
