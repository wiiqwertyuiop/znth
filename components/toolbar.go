package components

import (
	"encoding/json"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"znth/audio"
	"znth/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var stems []model.Stem

func PlayAction() {
	audio.TogglePlay()
}

func StopAction() {
	audio.Pause()
	audio.SetMusicPosition(0)
}

func LoadProjectFolder(folder string, mixers *fyne.Container) []model.Stem {
	files, err := os.ReadDir(folder)
	if err != nil {
		panic(err)
	}

	stems = stems[:0]
	savedData, _ := LoadStemData(folder)
	savedStemData := savedData.Stems

	savedMasterVolume, exists := savedStemData[0]
	masterVolume := audio.SliderToGain(30.0 / 100)
	if exists {
		masterVolume = savedMasterVolume.VolumeAdjust
		audio.SetMasterVolume(savedMasterVolume.VolumeAdjust)
	}
	// Master Volume
	volume := NewVerticalSlider(0, 100, audio.GainToSlider(masterVolume)*100)
	volume.OnChanged = func(v float64) {
		audio.SetMasterVolume(audio.SliderToGain(v / 100.0))
	}

	border := canvas.NewRectangle(color.NRGBA{
		R: 60,
		G: 60,
		B: 60,
		A: 255,
	})
	border.StrokeColor = color.Black
	border.StrokeWidth = 1

	bg := canvas.NewRectangle(color.NRGBA{
		R: 144,
		G: 0,
		B: 0,
		A: 144,
	})

	label := container.NewPadded(widget.NewLabel("Master"))

	nameLabel := container.NewMax(
		bg,
		label,
	)

	channel := container.NewStack(
		border,
		container.NewBorder(
			nil,
			nameLabel,
			nil,
			nil,
			volume,
		),
	)

	mixers.Add(channel)

	for id, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.ToLower(filepath.Ext(file.Name())) == ".wav" {
			fullPath := filepath.Join(folder, file.Name())
			println("Loading... ", fullPath)
			data, info, err := audio.LoadWavFloat32(fullPath)
			if err != nil {
				panic(err)
			}

			index := len(stems)

			savedVolume, exists := savedStemData[id]
			defaultVolume := audio.SliderToGain(50.0 / 100.0)
			if exists {
				defaultVolume = savedVolume.VolumeAdjust
			}

			volume := NewVerticalSlider(0, 100, audio.GainToSlider(defaultVolume)*100)
			stem := model.Stem{
				Id:           id,
				Data:         data,
				Info:         info,
				VolumeAdjust: defaultVolume,
			}

			volume.OnChanged = func(v float64) {
				stems[index].VolumeAdjust = audio.SliderToGain(v / 100.0)
			}

			border := canvas.NewRectangle(color.NRGBA{
				R: 60,
				G: 60,
				B: 60,
				A: 255,
			})
			border.StrokeColor = color.Black
			border.StrokeWidth = 1

			// TODO: Determine track color
			bg := canvas.NewRectangle(color.NRGBA{
				R: 144,
				G: 238,
				B: 144,
				A: 144,
			})

			label := container.NewPadded(widget.NewLabel(file.Name()))

			nameLabel := container.NewMax(
				bg,
				label,
			)

			channel := container.NewStack(
				border,
				container.NewBorder(
					nil,
					nameLabel,
					nil,
					nil,
					volume,
				),
			)

			mixers.Add(channel)

			stems = append(stems, stem)
		}
	}
	mixers.Refresh()
	audio.StartStream(stems)
	audio.Pause()
	return stems
}

func SaveStemData(path string) error {
	project := model.SongSave{
		Stems: make(map[int]model.SavedStemData),
	}

	for _, stem := range stems {
		project.Stems[stem.Id] = model.SavedStemData{
			VolumeAdjust: stem.VolumeAdjust,
		}
	}

	// Save master volume
	project.Stems[0] = model.SavedStemData{
		VolumeAdjust: audio.GetMasterVolume(),
	}

	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path+"/config.json", data, 0644)
}

func LoadStemData(path string) (model.SongSave, error) {
	var project model.SongSave

	filename := path + "/config.json"

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return project, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return project, err
	}

	err = json.Unmarshal(data, &project)
	if err != nil {
		return project, err
	}

	return project, nil
}
