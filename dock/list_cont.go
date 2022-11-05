// Copyright (C) 2022  Shanhu Tech Inc.
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

package dock

import (
	"net/url"
)

// ContListInfo is the container info got from listing.
type ContListInfo struct {
	ID      string `json:"Id"`
	Names   []string
	Image   string
	ImageID string
	Labels  map[string]string
}

// ListContsWithLabel lists all containers with the given label.
func ListContsWithLabel(c *Client, label string) ([]*ContListInfo, error) {
	filters, err := labelFilters(label)
	if err != nil {
		return nil, err
	}
	q := make(url.Values)
	q.Add("all", "true")
	q.Add("filters", filters)

	var infos []*ContListInfo
	if err := c.jsonGet("containers/json", q, &infos); err != nil {
		return nil, err
	}
	return infos, nil
}
