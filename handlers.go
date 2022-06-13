package main

import (
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/labstack/echo/v4"
	ua "github.com/wux1an/fake-useragent"
)

func MasterPlaylistHandler(c echo.Context) error {
	channel, err := ChannelFromContext(c)
	if err != nil {
		return c.String(http.StatusNotFound, "not found")
	}
	res, err := sendRequest(channel.URL.String(), nil)
	if err != nil {
		return err
	}
	data, err := parseM3u8FromResponse(res, nil)
	if err != nil {
		return err
	}
	return c.Stream(res.StatusCode, "application/x-mpegURL", data)
}

func MediaPlaylistHandler(c echo.Context) error {
	channelSlug := c.Param("channelSlug")
	if channelSlug == "" {
		return c.String(http.StatusNotFound, "not found")
	}
	if channelSlug == "" {
		return c.String(http.StatusNotFound, "not found")
	}
	channel, err := ChannelFromContext(c)
	if err != nil {
		return c.String(http.StatusNotFound, "not found")
	}
	cookies := cookiesFromQuery(c.QueryParams())
	res, err := sendRequest(channel.RenditionURL(c.Param("rendition")), cookies)
	if err != nil {
		return err
	}
	data, err := parseM3u8FromResponse(res, cookies)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return c.Stream(res.StatusCode, "application/x-mpegURL", data)
}

func SegmentHandler(c echo.Context) error {
	channel, err := ChannelFromContext(c)
	if err != nil {
		return c.String(http.StatusNotFound, "not found")
	}
	cookies := cookiesFromQuery(c.QueryParams())
	res, err := sendRequest(channel.SegmentURL(c.Param("rendition"), c.Param("segment")), cookies)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return c.Stream(res.StatusCode, "application/x-mpegURL", res.Body)
}

func cookiesFromQuery(qs url.Values) []*http.Cookie {
	domain := qs.Get("domain")
	cookies := []*http.Cookie{}
	for name, values := range qs {
		if len(values) == 0 {
			continue
		}
		if name == "domain" {
			continue
		}
		cookies = append(cookies, &http.Cookie{
			Name:    name,
			Value:   values[0],
			Path:    "/",
			Domain:  domain,
			Expires: time.Now().Add(24 * time.Hour),
		})
	}
	return cookies
}

func sendRequest(url string, cookies []*http.Cookie) (*http.Response, error) {
	req, err := retryablehttp.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}
	req.Header.Set("User-Agent", ua.RandomType(ua.Desktop))
	return httpClient.Do(req)
}
