package audio

import (
	"math"
	"znth/model"
	"znth/state"

	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"github.com/gordonklaus/portaudio"
)

var masterVolume float32 = SliderToGain(30.0 / 100.0)

var stream *portaudio.Stream = nil

var finishedChan = make(chan bool)

func Initialize() {
	portaudio.Initialize()
}

func StartStream(stems []*model.Stem, state *state.State) {

	var err error
	stemPeaks := make([]float32, len(stems))

	stream, err = portaudio.OpenDefaultStream(
		0,
		2,
		float64(stems[0].Info.SampleRate),
		512,
		func(out []float32) {

			position := state.Playback.Position.Load()
			clear(stemPeaks)

			if position >= int64(len(stems[0].Data)) {
				for i := range out {
					out[i] = 0
				}

				fyne.Do(func() { Stop(state) })
				return
			}

			for i := 0; i < len(out); i += 2 {

				var left float32
				var right float32

				samplePosition := position + int64(i)

				for index, stem := range stems {

					if samplePosition+1 < int64(len(stem.Data)) {

						stemLeft := float32(stem.Data[samplePosition]) / 32768.0 * stem.VolumeAdjust
						stemRight := float32(stem.Data[samplePosition+1]) / 32768.0 * stem.VolumeAdjust

						// Track this stem's peak
						peak := float32(math.Max(
							math.Abs(float64(stemLeft)),
							math.Abs(float64(stemRight)),
						))

						if peak > stemPeaks[index] {
							stemPeaks[index] = peak
						}

						// Mix
						left += stemLeft
						right += stemRight
					}
				}

				out[i] = float32(math.Tanh(float64(left * masterVolume)))
				out[i+1] = float32(math.Tanh(float64(right * masterVolume)))

				state.Playback.Position.Store(position + int64(len(out)))
			}

			// Store peaks once per buffer
			for i, stem := range stems {
				stem.Peak.Store(math.Float32bits(stemPeaks[i]))
			}
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}

func GetMasterVolume() float32 {
	return masterVolume
}

func SetMasterVolume(vol float32) {
	masterVolume = vol
}

func IsStreamActive() bool {
	return stream != nil
}

func Play(state *state.State) {
	if IsStreamActive() {
		state.PlaybackChange(model.PlaybackPlaying)
		stream.Start()
	}
}

func Pause(state *state.State) {
	if IsStreamActive() {
		state.PlaybackChange(model.PlaybackPaused)
		stream.Stop()
		for _, stem := range state.Project.Channels.Stems {
			stem.Peak.Store(math.Float32bits(0))
		}
	}
}

func Stop(state *state.State) {
	if IsStreamActive() {
		state.PlaybackChange(model.PlaybackStopped)
		state.Playback.Position.Store(0)
		stream.Stop()
		for _, stem := range state.Project.Channels.Stems {
			stem.Peak.Store(math.Float32bits(-1))
		}
	}
}

func TogglePlay(state *state.State) {
	if state.Playback.State != model.PlaybackPlaying {
		Play(state)
	} else {
		Pause(state)
	}
}

func LoadWavInt16(filename string) ([]int16, model.WavInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, model.WavInfo{}, err
	}
	defer file.Close()

	// RIFF header
	var riff [4]byte
	file.Read(riff[:])

	if string(riff[:]) != "RIFF" {
		return nil, model.WavInfo{}, fmt.Errorf("not a WAV file")
	}

	file.Seek(4, io.SeekCurrent) // skip file size

	var wave [4]byte
	file.Read(wave[:])

	if string(wave[:]) != "WAVE" {
		return nil, model.WavInfo{}, fmt.Errorf("not WAVE format")
	}

	var info model.WavInfo
	var audioFormat uint16
	var dataSize uint32

	// Read chunks
	for {
		var chunkID [4]byte
		err := binary.Read(file, binary.LittleEndian, &chunkID)
		if err != nil {
			return nil, info, err
		}

		var chunkSize uint32
		binary.Read(file, binary.LittleEndian, &chunkSize)

		switch string(chunkID[:]) {

		case "fmt ":
			binary.Read(file, binary.LittleEndian, &audioFormat)
			binary.Read(file, binary.LittleEndian, &info.Channels)

			binary.Read(file, binary.LittleEndian, &info.SampleRate)

			fmt.Println(
				"format:", audioFormat,
				"channels:", info.Channels,
				"rate:", info.SampleRate,
			)

			// Skip:
			// byte rate (4)
			// block align (2)
			// bits per sample (2)
			file.Seek(8, io.SeekCurrent)

			// skip any extra fmt bytes
			if chunkSize > 16 {
				file.Seek(int64(chunkSize-16), io.SeekCurrent)
			}

		case "data":
			dataSize = chunkSize

			goto READ_DATA

		default:
			// skip unknown chunks
			file.Seek(int64(chunkSize), io.SeekCurrent)
		}
	}

READ_DATA:

	if audioFormat != 1 {
		return nil, info, fmt.Errorf(
			"expected PCM WAV, got format %d",
			audioFormat,
		)
	}

	if dataSize%2 != 0 {
		return nil, info, fmt.Errorf("invalid int16 data size")
	}

	samples := make([]int16, dataSize/2)

	err = binary.Read(
		file,
		binary.LittleEndian,
		&samples,
	)

	if err != nil {
		return nil, info, err
	}

	return samples, info, nil
}

func KillStream(state *state.State) {
	if stream != nil {
		stream.Stop()
		stream.Close()
		stream = nil
	}

	state.Playback.Position.Store(0)
	state.Project.Channels.Stems = nil
	state.PlaybackChange(model.PlaybackStopped)
}

func Shutdown() {
	if IsStreamActive() {
		stream.Stop()
		stream.Close()
		stream = nil
	}
	portaudio.Terminate()
}
