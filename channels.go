package main

import (
	"bytes"
	_ "embed"

	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type Channel struct {
	Name string
	Slug string
	URL  *url.URL
}

func (c *Channel) UnmarshalJSON(data []byte) error {
	raw := struct {
		Name string
		Slug string
		URL  string
	}{}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	u, err := url.Parse(raw.URL)
	if err != nil {
		return err
	}
	*c = Channel{
		Name: raw.Name,
		Slug: raw.Slug,
		URL:  u,
	}
	return nil
}

func (c *Channel) SegmentURL(rendition string, segment string) string {
	return "https://" + c.URL.Host + filepath.Join(c.PlaylistBasePath(), rendition, segment)
}

func (c *Channel) RenditionURL(rendition string) string {
	return "https://" + c.URL.Host + filepath.Join(c.PlaylistBasePath(), rendition, "playlist.m3u8")
}

func (c *Channel) PlaylistBasePath() string {
	return filepath.Dir(c.URL.Path)
}

type Channels []*Channel

func (cs Channels) BySlug(slug string) (*Channel, bool) {
	for _, c := range cs {
		if c.Slug == slug {
			return c, true
		}
	}
	return nil, false
}

var (
	//go:embed channels.json
	channelsData []byte
)

func LoadChannels() (Channels, error) {
	var channels Channels
	err := json.NewDecoder(bytes.NewReader(channelsData)).Decode(&channels)
	return channels, err
}

const (
	channelContextKey = "channel"
)

func ChannelFromContext(c echo.Context) (*Channel, error) {
	channel, ok := c.Get(channelContextKey).(*Channel)
	if ok {
		return channel, nil
	}
	return nil, fmt.Errorf("channel not found")
}

func LoadChannelBySlug(channels Channels) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			channel, ok := channels.BySlug(c.Param("channelSlug"))
			if !ok {
				return c.String(http.StatusNotFound, "not found")
			}
			c.Set(channelContextKey, channel)
			return next(c)
		}
	}
}
