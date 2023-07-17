package downloader

import (
	"fmt"

	"github.com/gapidobri/librespot-golang/librespot/utils"
	"github.com/pkg/errors"
)

func (d *Downloader) DownloadPlaylist(playlistId string) error {
	playlist, err := d.session.Mercury().GetPlaylist(utils.Base62ToHex(playlistId))
	if err != nil {
		return errors.Wrap(err, "failed to get playlist")
	}

	playlist.GetLength()

	fmt.Println(playlist)

	// for _, item := range playlist.GetContents().Items {
	// 	fmt.Println(item.String())
	// }

	return nil
}
