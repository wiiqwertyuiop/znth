package components

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"znth/audio"
	"znth/model"
)

func loadProjectFolder(folder string) model.Channels {
	var stems []model.Stem
	var channels model.Channels

	files, err := os.ReadDir(folder)
	if err != nil {
		//panic(err) // todo
		return channels
	}

	savedData, _ := LoadStemData(folder)
	savedStemData := savedData.Stems

	savedMasterVolume, exists := savedStemData["0"]
	masterVolume := audio.SliderToGain(30.0 / 100)
	if exists {
		masterVolume = savedMasterVolume.VolumeAdjust
		audio.SetMasterVolume(savedMasterVolume.VolumeAdjust)
	}

	channels.MasterVolume = masterVolume

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.ToLower(filepath.Ext(file.Name())) == ".wav" {
			fullPath := filepath.Join(folder, file.Name())
			println("Loading... ", fullPath)
			data, info, err := audio.LoadWavFloat32(fullPath)
			if err != nil {
				//panic(err) // TODO
				continue
			}

			savedVolume, exists := savedStemData[file.Name()]
			defaultVolume := audio.SliderToGain(50.0 / 100.0)
			if exists {
				defaultVolume = savedVolume.VolumeAdjust
			}

			stem := model.Stem{
				Id:           file.Name(),
				Data:         data,
				Info:         info,
				VolumeAdjust: defaultVolume,
			}

			stems = append(stems, stem)
		}
	}

	channels.Stems = stems
	return channels
}

func SaveStemData(path string, stems []model.Stem) error {
	project := model.SongSave{
		Stems: make(map[string]model.SavedStemData),
	}

	for _, stem := range stems {
		project.Stems[stem.Id] = model.SavedStemData{
			VolumeAdjust: stem.VolumeAdjust,
		}
	}

	// Save master volume
	project.Stems["0"] = model.SavedStemData{
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
