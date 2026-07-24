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

	"github.com/gordonklaus/portaudio"
)

var position = 0
var masterVolume float32 = SliderToGain(30.0 / 100.0)

var stream *portaudio.Stream = nil

var finishedChan = make(chan bool)

func Initialize() {
	portaudio.Initialize()
}

func StartStream(stems []model.Stem, state *state.State) {

	var err error
	stream, err = portaudio.OpenDefaultStream(
		0, // input channels
		2, // stereo
		float64(stems[0].Info.SampleRate),
		512,
		func(out []float32) {

			// Check if playback finished
			if position >= len(stems[0].Data) {
				for i := range out {
					out[i] = 0
				}

				// Notify the goroutine
				select {
				case finishedChan <- true:
				default:
				}

				return
			}

			// Start at 1 so we skip the master channel
			for i := 0; i < len(out); i += 2 {

				var left float32
				var right float32

				for _, val := range stems {

					if position+1 < len(val.Data) {
						left += val.Data[position] * val.VolumeAdjust
						right += val.Data[position+1] * val.VolumeAdjust
					}
				}

				// master volume
				out[i] = float32(math.Tanh(float64(left * masterVolume)))
				out[i+1] = float32(math.Tanh(float64(right * masterVolume)))

				position += 2
			}
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	// Wait for playback to finish without blocking audio
	go func() {
		<-finishedChan

		KillStream(state)
	}()
}

func SetMusicPosition(pos int) {
	position = pos
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
	}
}

func Stop(state *state.State) {
	if IsStreamActive() {
		state.PlaybackChange(model.PlaybackStopped)
		SetMusicPosition(0)
		stream.Stop()
	}
}

func TogglePlay(state *state.State) {
	if state.Playback.State != model.PlaybackPlaying {
		Play(state)
	} else {
		Pause(state)
	}
}

func LoadWavFloat32(filename string) ([]float32, model.WavInfo, error) {
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

	if audioFormat != 3 {
		return nil, info, fmt.Errorf(
			"expected IEEE float WAV, got format %d",
			audioFormat,
		)
	}

	if dataSize%4 != 0 {
		return nil, info, fmt.Errorf("invalid float data size")
	}

	samples := make([]float32, dataSize/4)

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
		position = 0
		state.PlaybackChange(model.PlaybackStopped)
	}
}

func Shutdown() {
	stream.Stop()
	stream.Close()
	stream = nil
	position = 0
	portaudio.Terminate()
}
