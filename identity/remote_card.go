package identity

import (
	"context"
	"net/url"
	"sync"
	"time"

	"shanhu.io/g/errcode"
	"shanhu.io/g/httputil"
)

// RemoteCard is a remote identity.
type RemoteCard struct {
	client  *httputil.Client
	apiPath string

	mu          sync.Mutex
	cache       *Identity
	cacheExpire time.Time

	now func() time.Time
}

// NewRemoteCard creates a new remote card
func NewRemoteCard(u *url.URL) *RemoteCard {
	server := &url.URL{
		Scheme: u.Scheme,
		User:   u.User,
		Host:   u.Host,
	}
	apiPath := u.Path
	client := &httputil.Client{Server: server}

	return &RemoteCard{
		client:  client,
		apiPath: apiPath,
		now:     time.Now,
	}
}

// Refresh forces a refresh of the cached identity.
// The context is currently ignored.
func (c *RemoteCard) Refresh(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Lock()
	return c.refresh(ctx)
}

func (c *RemoteCard) refresh(ctx context.Context) error {
	// Must already holding mutext.
	req := new(GetIDRequest)
	id := new(Identity)
	if err := c.client.Call(c.apiPath, req, id); err != nil {
		return err
	}
	c.cache = id

	// Refresh at least once per hour.
	c.cacheExpire = c.now().Add(time.Hour)
	return nil
}

func (c *RemoteCard) ensure(ctx context.Context) error {
	if c.cache == nil {
		return c.refresh(ctx)
	}
	now := c.now()
	if now.After(c.cacheExpire) {
		return c.refresh(ctx)
	}

	const grace = time.Minute * 10
	cut := now.Add(grace).Unix()
	for _, k := range c.cache.PublicKeys {
		if k.NotValidAfter >= cut {
			return nil // Found a key that isn't expiring soon.
		}
	}

	return c.refresh(ctx)
}

// Identity returns the identity that is fetched from the remote API
// endpoint.
func (c *RemoteCard) Identity(ctx context.Context) (*Identity, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.ensure(ctx); err != nil {
		return nil, err
	}
	if c.cache == nil {
		return nil, errcode.NotFoundf("identity not found")
	}
	return c.cache, nil
}
