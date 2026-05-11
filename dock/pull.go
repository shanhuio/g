package dock

import (
	"net/url"
)

// PullImage pulls the specified image and save it as the related tag.
func PullImage(c *Client, image, tag string) error {
	q := make(url.Values)
	q.Add("fromImage", image)
	if tag != "" {
		q.Add("tag", tag)
	}

	sink := newStreamSink()
	if err := c.post("images/create", q, nil, sink); err != nil {
		return err
	}
	return sink.waitDone()
}
