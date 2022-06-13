package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/grafov/m3u8"
)

func masterPlaylistWithCookieAuth(master *m3u8.MasterPlaylist, qs url.Values) (*bytes.Buffer, error) {
	for i := range master.Variants {
		u, err := url.Parse(master.Variants[i].URI)
		if err != nil {
			return nil, err
		}
		u.RawQuery = qs.Encode()
		master.Variants[i].URI = u.String()
	}
	return master.Encode(), nil
}

func mediaPlaylistWithCookieAuth(media *m3u8.MediaPlaylist, qs url.Values) (*bytes.Buffer, error) {
	for i := range media.Segments {
		if media.Segments[i] == nil {
			continue
		}
		u, err := url.Parse(media.Segments[i].URI)
		if err != nil {
			return nil, err
		}
		u.RawQuery = qs.Encode()
		media.Segments[i].URI = u.String()
	}
	return media.Encode(), nil
}

func parseM3u8FromResponse(res *http.Response, cookies []*http.Cookie) (*bytes.Buffer, error) {
	defer res.Body.Close()
	playlist, _, err := m3u8.DecodeFrom(res.Body, false)
	if err != nil {
		return nil, err
	}
	if cookies == nil {
		cookies = res.Cookies()
	}
	return addAuthQueryParams(playlist, cookies)
}

func addAuthQueryParams(playlist m3u8.Playlist, cookies []*http.Cookie) (*bytes.Buffer, error) {
	qs := url.Values{}
	for _, cookie := range cookies {
		qs.Set(cookie.Name, cookie.Value)
		if cookie.Domain != "" {
			qs.Set("domain", cookie.Domain)
		}
	}
	if master, ok := playlist.(*m3u8.MasterPlaylist); ok {
		return masterPlaylistWithCookieAuth(master, qs)
	}
	media, ok := playlist.(*m3u8.MediaPlaylist)
	if !ok {
		return nil, fmt.Errorf("unknown playlist type: %s", playlist)
	}
	return mediaPlaylistWithCookieAuth(media, qs)
}
