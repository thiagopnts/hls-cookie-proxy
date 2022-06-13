package main

import (
	"log"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/labstack/echo/v4"
)

var (
	retryWait  = 0 * time.Second
	timeout    = 5 * time.Second
	httpClient = newHTTPClient()
)

func newHTTPClient() *retryablehttp.Client {
	httpClient := retryablehttp.NewClient()
	httpClient.HTTPClient.Timeout = timeout
	httpClient.RetryMax = 2
	httpClient.RetryWaitMin = retryWait
	httpClient.RetryWaitMax = retryWait
	return httpClient
}

func main() {
	channels, err := LoadChannels()
	if err != nil {
		log.Fatalf("error loading channels.json file: %s", err)
	}
	e := echo.New()
	e.Use(LoadChannelBySlug(channels))
	e.GET("/:channelSlug/playlist.m3u8", MasterPlaylistHandler)
	e.GET("/:channelSlug/:rendition/playlist.m3u8", MediaPlaylistHandler)
	e.GET("/:channelSlug/:rendition/:segment", SegmentHandler)
	e.Logger.Fatal(e.Start(":1323"))
}
