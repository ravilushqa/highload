package centrifugoclient

import (
	"github.com/centrifugal/gocent"
)

func New(addr, key string) *gocent.Client {
	return gocent.New(gocent.Config{
		Addr:       addr,
		Key:        key,
		HTTPClient: nil,
	})
}
