package you

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aaaasmile/live-omxctrl/db"
)

type YoutubePl struct {
	YoutubeInfo *InfoLink
	URI         string
	TmpInfo     string
}

func (yp *YoutubePl) GetStatusSleepTime() int {
	return 1700
}

func (yp *YoutubePl) GetURI() string {
	return yp.URI
}

func (yp *YoutubePl) GetTitle() string {
	if yp.YoutubeInfo != nil {
		return yp.YoutubeInfo.Title
	}
	return ""
}

func (yp *YoutubePl) Name() string {
	return "youtube"
}

func (yp *YoutubePl) CheckStatus(chHistoryItem chan *db.HistoryItem) (bool, error) {
	if yp.YoutubeInfo == nil {
		info, err := readLinkDescription(yp.URI, yp.TmpInfo)
		yp.YoutubeInfo = info
		if err != nil {
			return false, err
		}

		hi := db.HistoryItem{
			URI:           yp.URI,
			Title:         info.Title,
			Description:   info.Description,
			DurationInSec: info.Duration,
			Type:          yp.Name(),
			Duration:      time.Duration(int64(info.Duration) * int64(time.Second)).String(),
		}
		chHistoryItem <- &hi
	}
	return true, nil
}

func (yp *YoutubePl) GetDescription() string {
	if yp.YoutubeInfo != nil {
		return yp.YoutubeInfo.Description
	}
	return ""
}

func (yp *YoutubePl) IsUriForMe(uri string) bool {
	if strings.Contains(uri, "you") && strings.Contains(uri, "https") {
		log.Println("this is youtube URL ", uri)
		yp.URI = uri
		return true
	}
	return false
}

func (yp *YoutubePl) GetStreamerCmd(cmdLineArr []string) string {
	args := strings.Join(cmdLineArr, " ")
	cmd := fmt.Sprintf("omxplayer %s `%s -f mp4 -g %s`", args, getYoutubePlayer(), yp.URI)
	return cmd
}

func getYoutubePlayer() string {
	return "you" + "tube" + "-" + "dl"
}