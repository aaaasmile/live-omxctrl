package playlist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	dirPlaylistData = "playlist-data"
)

const (
	lastPlayedInfo = "info.json"
)

type PlayinfoLast struct {
	Playlist string
	URI      string
}

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
	Name    string
	Created string
	List    []*PlayItem
}

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

type LLPlayItem struct {
	PlayItem *PlayItem
	Next     *LLPlayItem
	Previous *LLPlayItem
}

func NewLLPlayItem(nx, pr *LLPlayItem, plit *PlayItem) *LLPlayItem {
	res := LLPlayItem{
		PlayItem: plit,
		Next:     nx,
		Previous: pr,
	}
	return &res
}

type LLPlayList struct {
	Name      string
	Count     int
	FirstItem *LLPlayItem
	LastItem  *LLPlayItem
	CurrItem  *LLPlayItem
}

func (ll *LLPlayList) First() {
	ll.CurrItem = ll.FirstItem
}

func (ll *LLPlayList) Last() {
	ll.CurrItem = ll.LastItem
}

func (ll *LLPlayList) Next() (*PlayItem, bool) {
	ll.CurrItem = ll.CurrItem.Next
	if ll.CurrItem == nil {
		ll.CurrItem = ll.FirstItem
	}
	if ll.CurrItem != nil {
		return ll.CurrItem.PlayItem, ll.CurrItem.PlayItem != nil
	} else {
		return nil, false
	}
}

func (ll *LLPlayList) Previous() (*PlayItem, bool) {
	ll.CurrItem = ll.CurrItem.Previous
	if ll.CurrItem == nil {
		ll.CurrItem = ll.LastItem
	}
	if ll.CurrItem != nil {
		return ll.CurrItem.PlayItem, ll.CurrItem.PlayItem != nil
	} else {
		return nil, false
	}
}

func (ll *LLPlayList) CheckCurrent() (*PlayItem, bool) {
	if ll.FirstItem == nil ||
		ll.LastItem == nil ||
		ll.CurrItem == nil {
		return nil, false
	}
	if ll.CurrItem.PlayItem == nil {
		return nil, false
	}
	return ll.CurrItem.PlayItem, true
}
func (ll *LLPlayList) IsEmpty() bool {
	return ll.FirstItem == nil ||
		ll.LastItem == nil
}

// func (ll *LLPlayList) StartPlay() error {
// 	// TODO
// 	return nil
// }

func GetCurrentPlaylist() (*LLPlayList, error) {
	res := &LLPlayList{}
	var err error
	infopath := filepath.Join(dirPlaylistData, lastPlayedInfo)
	if _, err = os.Stat(infopath); err == nil {
		raw, err := ioutil.ReadFile(infopath)
		if err != nil {
			return res, err
		}
		pllast := PlayinfoLast{}
		if err := json.Unmarshal(raw, &pllast); err != nil {
			return res, err
		}
		log.Println("Last played info ", pllast)
		return res, fmt.Errorf("TODO provides the last playlist")
	}

	fileitems, err := ioutil.ReadDir(dirPlaylistData)
	if err != nil {
		return res, err
	}
	plname := ""
	for _, item := range fileitems {
		if !item.IsDir() {
			log.Println("Selected list ", item)
		}
		nn := item.Name()
		plname = nn
		if strings.Contains(nn, "default") {
			break
		}
	}
	if plname != "" {
		res, err = BuildLListFromfile(filepath.Join(dirPlaylistData, plname))
		if err != nil {
			return res, err
		}
	}

	if res.IsEmpty() {
		return res, fmt.Errorf("Nothing to play")
	}
	log.Println("playlist len", res.Count)
	return res, nil
}

func BuildLListFromfile(fname string) (*LLPlayList, error) {
	log.Println("Building LL list from file ", fname)
	res := &LLPlayList{}
	raw, err := ioutil.ReadFile(fname)
	if err != nil {
		return res, err
	}
	list := PlayList{}
	if err := json.Unmarshal(raw, &list); err != nil {
		return res, err
	}
	log.Println("Playlist is ", list)
	res.Name = list.Name
	var llitem, pr *LLPlayItem
	for _, item := range list.List {
		pr = llitem
		llitem = NewLLPlayItem(nil, pr, item)
		if pr != nil {
			pr.Next = llitem
		}
		if res.CurrItem == nil {
			res.CurrItem = llitem
			res.FirstItem = llitem
		}
		res.LastItem = llitem
		res.Count++
	}

	return res, nil
}

func CheckIfPlaylistExist(plname string) error {
	path := filepath.Join(dirPlaylistData, plname)
	_, err := os.Stat(path)
	return err
}
