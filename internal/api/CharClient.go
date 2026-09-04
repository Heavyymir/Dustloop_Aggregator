package api

import (
	"net/http"
	"time"

	"github.com/Heavyymir/Dustloop_Aggregator/internal/CharDataCache"
)

// Temporary Client Struct, will move to internal API package later
type Client struct {
	HttpClient http.Client
	Cache      *CharDataCache.Cache
}

func NewClient() *Client {
	return &Client{
		HttpClient: http.Client{
			Timeout: 15 * time.Second,
		},
		Cache:      CharDataCache.NewCache(5 * time.Second),
	}
}
