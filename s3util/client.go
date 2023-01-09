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

package s3util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	minio "github.com/minio/minio-go/v7"
	creds "github.com/minio/minio-go/v7/pkg/credentials"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/httputil"
)

func makeMinioClient(conf *Config, cred *Credential) (*minio.Client, error) {
	opt := &minio.Options{
		Creds:  creds.NewStaticV4(cred.Key, cred.Secret, ""),
		Secure: true,
	}
	return minio.New(conf.Endpoint, opt)
}

func minioError(err error) error {
	if err == nil {
		return nil
	}

	if status := minio.ToErrorResponse(err).StatusCode; status != 0 {
		return httputil.AddErrCode(status, err)
	}
	return err
}

// Client gives a client for accessing an S3-compatible storage service.
type Client struct {
	client   *minio.Client
	config   *Config
	bucket   string
	basePath string
}

// NewClient creates a new client for accessing a storage endpoint with a base
// path prefix.
func NewClient(config *Config, cred *Credential) (*Client, error) {
	client, err := makeMinioClient(config, cred)
	if err != nil {
		return nil, err
	}
	return &Client{
		client:   client,
		config:   config,
		bucket:   config.Bucket,
		basePath: config.BasePath,
	}, nil
}

// BasePath returns the base path of the client.
func (c *Client) BasePath() string { return c.basePath }

// BaseURL returns the S3 URL of the base path.
func (c *Client) BaseURL() *url.URL {
	return &url.URL{
		Scheme: "s3",
		Host:   c.bucket,
		Path:   c.basePath,
	}
}

func (c *Client) path(p string) string {
	if c.basePath == "" {
		return p
	}
	return path.Join(c.basePath, p)
}

// PresignGet presigns a URL for a future GET. The URL expires in d.
func (c *Client) PresignGet(ctx C, p string, exp time.Duration) (
	*url.URL, error,
) {
	var v url.Values
	return c.client.PresignedGetObject(ctx, c.bucket, c.path(p), exp, v)
}

// PresignPut presigns a URL for a future PUT. The URL expires in d.
func (c *Client) PresignPut(ctx C, p string, exp time.Duration) (
	*url.URL, error,
) {
	return c.client.PresignedPutObject(ctx, c.bucket, c.path(p), exp)
}

// GetBytes gets an object.
func (c *Client) GetBytes(ctx C, p string) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := c.GetInto(ctx, p, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GetInto gets an object and writes its content into w.
func (c *Client) GetInto(ctx C, p string, w io.Writer) error {
	getError := func(err error) error {
		return errcode.Annotatef(minioError(err), "get %q", p)
	}

	p = c.path(p)
	var opts minio.GetObjectOptions
	obj, err := c.client.GetObject(ctx, c.bucket, p, opts)
	if err != nil {
		return getError(err)
	}
	defer obj.Close()
	if _, err := io.Copy(w, obj); err != nil {
		return getError(err)
	}
	return nil
}

// GetJSON gets a JSON object.
func (c *Client) GetJSON(ctx C, p string, v interface{}) error {
	bs, err := c.GetBytes(ctx, p)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}

// Put puts an object into the bucket.
func (c *Client) Put(
	ctx C, p string, r io.Reader, n int64, contentType string,
) error {
	p = c.path(p)
	var opts minio.PutObjectOptions
	if contentType != "" {
		opts.ContentType = contentType
	}
	if _, err := c.client.PutObject(ctx, c.bucket, p, r, n, opts); err != nil {
		return errcode.Annotatef(minioError(err), "put %q", p)
	}
	return nil
}

// Copy copies an object.
func (c *Client) Copy(ctx C, from, to string) error {
	dest := minio.CopyDestOptions{
		Bucket: c.bucket,
		Object: c.path(to),
	}
	src := minio.CopySrcOptions{
		Bucket: c.bucket,
		Object: c.path(from),
	}
	_, err := c.client.CopyObject(ctx, dest, src)
	return err
}

// PutFile copies a file.
func (c *Client) PutFile(ctx C, p, fp string) error {
	f, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}
	return c.Put(ctx, p, f, stat.Size(), "")
}

// PutBytes saves an object.
func (c *Client) PutBytes(ctx C, p string, data []byte) error {
	r := bytes.NewReader(data)
	typ := http.DetectContentType(data)
	return c.Put(ctx, p, r, int64(len(data)), typ)
}

// PutJSON puts a JSON object.
func (c *Client) PutJSON(ctx C, p string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.PutBytes(ctx, p, bs)
}

// Delete deletes an object.
func (c *Client) Delete(ctx C, p string) error {
	var opt minio.RemoveObjectOptions
	if err := c.client.RemoveObject(
		ctx, c.bucket, c.path(p), opt,
	); err != nil {
		return errcode.Annotatef(minioError(err), "delete %q", p)
	}
	return nil
}

// Stat returns the object info a particular object.
func (c *Client) Stat(ctx C, p string) (*minio.ObjectInfo, error) {
	var opts minio.StatObjectOptions
	info, err := c.client.StatObject(ctx, c.bucket, c.path(p), opts)
	if err != nil {
		return nil, errcode.Annotatef(minioError(err), "stat %q", p)
	}
	return &info, nil
}

// ListAll lists all objects in the bucket.
func (c *Client) ListAll(ctx C) ([]*minio.ObjectInfo, error) {
	var objs []*minio.ObjectInfo
	opt := minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    c.path(""),
	}
	ch := c.client.ListObjects(ctx, c.bucket, opt)
	for obj := range ch {
		if obj.Err != nil {
			return nil, obj.Err
		}
		cp := obj
		objs = append(objs, &cp)
	}
	return objs, nil
}
