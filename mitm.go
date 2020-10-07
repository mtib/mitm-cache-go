package mitm

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type (
	// Cache can make requests to a known mitm-cache service. If the cached result for a given url is older than `maxage` a new request will be made and the cache updated.
	Cache interface {
		// Request requests the url, only accept data which is at most `maxage` old.
		Request(url string, maxage time.Duration) (io.ReadCloser, error)
		// RequestNew also requests the url, but invalidates the cache, forcing the mitm server to make a new request.
		RequestNew(url string) (io.ReadCloser, error)
	}

	innerCache struct {
		upstream string
		key      string
	}
)

// New creates a new Cache. Upstream has to be the URI of a mitm-cache server.
func New(upstream string, key string) Cache {
	return innerCache{
		upstream: upstream,
		key:      key,
	}
}

func (c innerCache) Request(url string, maxage time.Duration) (io.ReadCloser, error) {
	if maxage < 0 {
		return nil, errors.New("Negative maxage not allowed, use 0 for invalidation")
	}
	resp, err := http.Get(fmt.Sprintf("%s/request/%d/%s/%s", c.upstream, maxage, base64.URLEncoding.EncodeToString([]byte(url)), c.key))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c innerCache) RequestNew(url string) (io.ReadCloser, error) {
	return c.Request(url, 0)
}
