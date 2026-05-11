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
