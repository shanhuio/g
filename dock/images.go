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

package dock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"

	"shanhu.io/g/errcode"
)

// RemoveImageOptions defines options when docker images are removed.
type RemoveImageOptions struct {
	Force   bool // Force remove the image
	NoPrune bool // NoPrune keeps un-tagged images after removal
}

// RemoveImage removes a local docker image.
func RemoveImage(c *Client, name string, opts *RemoveImageOptions) error {
	q := make(url.Values)
	if opts.Force {
		q.Add("force", "true")
	}
	if opts.NoPrune {
		q.Add("noprune", "true")
	}

	return c.del(path.Join("images", name), q)
}

// SaveImages saves built images as a tarball stream
// into the writer.
func SaveImages(c *Client, names []string, w io.Writer) error {
	v := make(url.Values)
	for _, name := range names {
		v.Add("names", name)
	}

	resp, err := c.get("images/get", v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	return err
}

// SaveImagesGz saves images to a gzipped tarball file.
func SaveImagesGz(c *Client, names []string, f string) error {
	gz, err := gzipCreate(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	if err := SaveImages(c, names, gz); err != nil {
		return err
	}
	return gz.Close()
}

// SaveImageGz saves a single image to a gzipped tarball file.
func SaveImageGz(c *Client, name, file string) error {
	return SaveImagesGz(c, []string{name}, file)
}

// LoadImages loads a tarball stream into Docker repository.
func LoadImages(c *Client, r io.Reader) error {
	sink := newStreamSink()
	if err := c.post("images/load", make(url.Values), r, sink); err != nil {
		return err
	}
	return sink.waitDone()
}

// LoadImagesFromFile loads a tarball file into Docker repository.
func LoadImagesFromFile(c *Client, f string) error {
	r, err := os.Open(f)
	if err != nil {
		return err
	}
	defer r.Close()

	return LoadImages(c, r)
}

// ImageInfo is the inspection result of an image.
type ImageInfo struct {
	ID          string `json:"Id,omitempty"`
	Parent      string
	VirtualSize int64
	RepoTags    []string `json:",omitempty"`
	RepoDigests []string `json:",omitempty"`
}

// InspectImage inspects a particular image.
func InspectImage(c *Client, name string) (*ImageInfo, error) {
	info := new(ImageInfo)
	p := path.Join("images", name, "json")
	if err := c.jsonGet(p, nil, info); err != nil {
		return nil, err
	}
	return info, nil
}

// HasImage checks if a particular image exists.
func HasImage(c *Client, name string) (bool, error) {
	if _, err := InspectImage(c, name); err != nil {
		if errcode.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// TagImage tags an image.
func TagImage(c *Client, name, repo, tag string) error {
	if tag == "" {
		return errcode.InvalidArgf("tag is empty")
	}
	if repo == "" {
		return errcode.InvalidArgf("repo is empty")
	}
	q := make(url.Values)
	q.Add("repo", repo)
	q.Add("tag", tag)

	p := path.Join("images", name, "tag")
	return c.poke(p, q)
}

// PruneImagesOption provides options for unused image pruning.
type PruneImagesOption struct {
	Unused bool
}

// PruneImages prunes the unused images.
func PruneImages(c *Client, opt *PruneImagesOption) error {
	filters := make(map[string][]string)
	filters["dangling"] = []string{fmt.Sprint(!opt.Unused)}
	bs, err := json.Marshal(filters)
	if err != nil {
		return errcode.InvalidArgf("marshal filter")
	}
	q := make(url.Values)
	q.Add("filters", string(bs))

	var resp struct{}
	return c.call("/images/prune", q, nil, &resp)
}

// ImageListInfo is the information of an image listing.
type ImageListInfo struct {
	ID       string `json:"Id,omitempty"`
	RepoTags []string
	Labels   map[string]string
}

// ListImages lists all images
func ListImages(c *Client) ([]*ImageListInfo, error) {
	var images []*ImageListInfo
	if err := c.jsonGet("/images/json", nil, &images); err != nil {
		return nil, err
	}
	return images, nil
}
