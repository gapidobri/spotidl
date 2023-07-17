package cmd

import (
	"fmt"
	"os"

	"github.com/gapidobri/spotidl/internal/downloader"
	"github.com/gapidobri/spotidl/internal/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "spotidl",
	Short: "Spotify music downloader",
	Long:  `A tool for directly downloading music from Spotify.`,
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		format := cmd.Flag("format").Value.String()

		dl, err := downloader.New(viper.GetString("username"), viper.GetString("password"))
		if err != nil {
			fmt.Println(err)
			return
		}

		urlType, id := downloader.GetTypeAndId(url)

		switch urlType {
		case downloader.UrlTypeTrack:
			err = dl.DownloadTrack(id, format)
		// case downloader.UrlTypeAlbum:
		// 	downloader.DownloadAlbum(id)
		// case downloader.UrlTypePlaylist:
		// 	err = dl.DownloadPlaylist(id)
		default:
			fmt.Println("not implemented")
		}

		if err != nil {
			fmt.Println(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	config.Init()

	rootCmd.Flags().StringP("format", "f", "mp3", "format to download the track in")
}
