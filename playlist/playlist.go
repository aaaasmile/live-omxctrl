package playlist

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const dirPlaylistData = "../playlist-data"

type PlayItemType int

const (
	PITMp3 = iota
	PITYoutube
	PITRadio
)

func (pi *PlayItemType) String() string {
	switch *pi {
	case PITMp3:
		return "Mp3"
	case PITYoutube:
		return "Youtube"
	case PITRadio:
		return "Radio"
	}
	return ""
}

type PlayItem struct {
	URI      string
	Info     string
	ItemType PlayItemType
}

type PlayList struct {
	Name       string
	LastPlayed string
	List       []*PlayItem
}

type LLPlayItem struct {
	PlayItem *PlayItem
	Next     *PlayItem
	Previous *PlayItem
}

type LLPlayList []*LLPlayItem

func (pl *PlayList) SavePlaylist(plname string) error {
	path := filepath.Join(dirPlaylistData, plname)
	log.Printf("Saving playlist file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Unable to save: %v", err)
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(pl)
}

func CheckIfPlaylistExist(plname string) error {
	path := filepath.Join(dirPlaylistData, plname)
	_, err := os.Stat(path)
	return err
}
