package downloader

import "regexp"

type UrlType string

const (
	UrlTypeTrack    UrlType = "track"
	UrlTypeAlbum    UrlType = "album"
	UrlTypePlaylist UrlType = "playlist"
)

var (
	trackRegex    = regexp.MustCompile("https://open.spotify.com/track/([[:alnum:]]{22})")
	albumRegex    = regexp.MustCompile("https://open.spotify.com/album/([[:alnum:]]{22})")
	playlistRegex = regexp.MustCompile("https://open.spotify.com/playlist/([[:alnum:]]{22})")
)

func GetTypeAndId(url string) (UrlType, string) {
	switch {
	case trackRegex.MatchString(url):
		return UrlTypeTrack, trackRegex.FindStringSubmatch(url)[1]
	case albumRegex.MatchString(url):
		return UrlTypeAlbum, albumRegex.FindStringSubmatch(url)[1]
	case playlistRegex.MatchString(url):
		return UrlTypePlaylist, playlistRegex.FindStringSubmatch(url)[1]
	}

	return "", ""
}
