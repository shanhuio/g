package states

import (
	"net/url"

	"shanhu.io/g/s3util"
)

type s3Back struct {
	client *s3util.Client
}

func newS3Back(client *s3util.Client) *s3Back {
	return &s3Back{client: client}
}

func (b *s3Back) Get(ctx C, key string) ([]byte, error) {
	return b.client.GetBytes(ctx, key)
}

func (b *s3Back) Put(ctx C, key string, data []byte) error {
	return b.client.PutBytes(ctx, key, data)
}

func (b *s3Back) Del(ctx C, key string) error {
	return b.client.Delete(ctx, key)
}

func (b *s3Back) URL() *url.URL {
	return b.client.BaseURL()
}
