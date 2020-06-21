package playlist

import (
	"testing"
	"time"
)

func TestCreateDefault(t *testing.T) {
	pli := PlayItem{
		URI:      "http://stream.srg-ssr.ch/m/rsc_de/aacp_96",
		Info:     "Radio Swiss Classic",
		ItemType: PITRadio,
	}
	list := make([]*PlayItem, 0)
	list = append(list, &pli)

	strl := PlayList{
		Name:    "RadioCH",
		List:    list,
		Created: time.Now().Format("02.01.2006 15:04:05"),
	}
	playlistName := "default"
	strl.SavePlaylist(playlistName)

	if err := CheckIfPlaylistExist("default"); err != nil {

		t.Error("Play list not created", playlistName, err)
	}
}
