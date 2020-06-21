package playlist

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	dirAbs := "../playlist-data"
	inifiles, err := ioutil.ReadDir(dirAbs)
	if err != nil {
		t.Error(err)
	}

	pli := PlayItem{
		URI:      "http://stream.srg-ssr.ch/m/rsc_de/aacp_96",
		Info:     "Radio Swiss Classic",
		ItemType: PITRadio,
	}
	list := make([]*PlayItem, 0)
	list = append(list, &pli)

	strl := PlayList{
		Name:       "RadioCH",
		List:       list,
		LastPlayed: time.Now().Format("02.01.2006 15:04:05"),
	}
	strl.SavePlaylist(filepath.Join(dirAbs, "default"))

	files, err := ioutil.ReadDir(dirAbs)
	if err != nil {
		t.Error(err)
	}
	if len(files) == len(inifiles) {
		t.Errorf("Playlist not created")
	}
}
