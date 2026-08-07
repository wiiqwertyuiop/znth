package components

import (
	"encoding/json"
	"errors"
	"io"
	"znth/model"
	"znth/state"

	"fyne.io/fyne/v2"
)

// TODO ERROR HANDLING, CHECK NIL, ETC.
func SaveSetlist(writer fyne.URIWriteCloser, state *state.State) error {
	if writer == nil {
		return errors.New("Folder selection cancelled")
	}

	defer writer.Close() // Ensure the file is closed to finalize the write

	setlistData := model.Setlist{Data: state.Project.SongNames}

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

func OpenSetlist(reader fyne.URIReadCloser, state *state.State) error {
	if reader == nil {
		return errors.New("Action canceled")
	}

	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	var setlist model.Setlist
	err = json.Unmarshal(data, &setlist)
	if err != nil {
		return err
	}

	state.Project.SongNames = setlist.Data
	LoadSong(setlist.Data[0].Location, state)
	return nil
}
