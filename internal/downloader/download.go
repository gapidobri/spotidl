package downloader

import (
	"encoding/binary"
	"fmt"

	spotify "github.com/gapidobri/librespot-golang/Spotify"
	"github.com/gapidobri/librespot-golang/librespot/core"
	"github.com/gapidobri/librespot-golang/librespot/utils"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/xlab/vorbis-go/decoder"
)

const sampleCount = 2048

type Downloader struct {
	session *core.Session
}

func New(username, password string) (*Downloader, error) {
	session, err := core.Login(username, password, "spotidl")
	if err != nil {
		return nil, err
	}

	return &Downloader{
		session: session,
	}, nil
}

func (d *Downloader) DownloadTrack(trackId string, format string) error {
	// Get track from id
	track, err := d.session.Mercury().GetTrack(utils.Base62ToHex(trackId))
	if err != nil {
		return errors.Wrap(err, "failed to get track")
	}

	fileName := fmt.Sprintf("%s - %s.%s", track.GetArtist()[0].GetName(), track.GetName(), format)

	fmt.Printf("Downloading %s\n", fileName)

	// Find the file with the highest bitrate
	trackFile, found := lo.Find(track.GetFile(), func(file *spotify.AudioFile) bool {
		return file.GetFormat() == spotify.AudioFile_OGG_VORBIS_320
	})
	if !found {
		return errors.Wrap(err, "audio file not found")
	}

	// Load the track into the player
	file, err := d.session.Player().LoadTrack(trackFile, track.GetGid())
	if err != nil {
		return errors.Wrap(err, "failed to load track")
	}

	// Create a ogg decoder for the track
	dec, err := decoder.New(file, sampleCount)
	if err != nil {
		return errors.Wrap(err, "failed to create decoder")
	}

	decoderInfo := dec.Info()

	channelCount := int(decoderInfo.Channels)
	sampleRate := decoderInfo.SampleRate

	// Start decoding the track
	go func() {
		dec.Decode()
		dec.Close()
	}()

	// Compile ffmpeg command
	ffmpegCmd := ffmpeg.
		Input("pipe:", ffmpeg.KwArgs{
			"f":  "f32le",
			"ac": channelCount,
			"ar": sampleRate,
		}).
		Output(fileName).
		Compile()

	// Get stdin pipe from ffmpeg
	ffmpegStdin, err := ffmpegCmd.StdinPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stdin pipe from ffmpeg")
	}

	out := make([]float32, sampleCount*channelCount)

	// Start writing to ffmpeg stdin
	go func() {
		defer ffmpegStdin.Close()

		for frame := range dec.SamplesOut() {
			if len(frame) > int(sampleCount) {
				frame = frame[:sampleCount]
			}

			var idx int
			for _, sample := range frame {
				if len(sample) > channelCount {
					sample = sample[:channelCount]
				}
				for i := range sample {
					out[idx] = sample[i]
					idx++
				}
			}

			binary.Write(ffmpegStdin, binary.LittleEndian, out)
		}
	}()

	// Run ffmpeg
	if err := ffmpegCmd.Run(); err != nil {
		return errors.Wrap(err, "failed to run ffmpeg")
	}

	return nil
}
