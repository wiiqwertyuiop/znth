package components

import (
	"encoding/json"
	"errors"
	"io"
	"slices"
	"znth/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Setlist data
var songNames []model.SongDetails

func CreateSetlist() *widget.List {
	return widget.NewList(
		func() int {
			return len(songNames)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			label.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
			return label
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(songNames[i].Name)
		},
	)
}

func AddToSetlist(reader fyne.ListableURI) {
	song := model.SongDetails{
		Name:     reader.Name(),
		Location: reader.Path(),
	}
	songNames = slices.Insert(songNames, 0, song)
}

// TODO ERROR HANDLING, CHECK NIL, ETC.
func SaveSetlist(writer fyne.URIWriteCloser) error {
	if writer == nil {
		return errors.New("Folder selection cancelled")
	}

	defer writer.Close() // Ensure the file is closed to finalize the write

	setlistData := model.Setlist{Data: songNames}

	// Example: Marshaling a Setlist slice to JSON
	data, err := json.Marshal(setlistData)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func OpenSetlist(reader fyne.URIReadCloser) {
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	var setlist model.Setlist
	err = json.Unmarshal(data, &setlist)
	if err != nil {
		panic(err)
	}

	songNames = setlist.Data
}

func GetSong(id int) model.SongDetails {
	return songNames[id]
}
