package live

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aaaasmile/live-omxctrl/conf"
	"github.com/aaaasmile/live-omxctrl/web/idl"
)

type PageCtx struct {
	RootUrl string
	Buildnr string
}

func getURLForRoute(uri string) string {
	arr := strings.Split(uri, "/")
	//fmt.Println("split: ", arr, len(arr))
	for i := len(arr) - 1; i >= 0; i-- {
		ss := arr[i]
		if ss != "" {
			if !strings.HasPrefix(ss, "?") {
				//fmt.Printf("Url for route is %s\n", ss)
				return ss
			}
		}
	}
	return uri
}

func APiHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	log.Println("Request: ", req.RequestURI)
	var err error
	switch req.Method {
	case "GET":
		err = handleGet(w, req)
	case "POST":
		log.Println("POST on ", req.RequestURI)
		err = handlePost(w, req)
	}
	if err != nil {
		log.Println("Error exec: ", err)
		http.Error(w, fmt.Sprintf("Internal error on execute: %v", err), http.StatusInternalServerError)
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.Printf("Service %s total call duration: %v\n", idl.Appname, elapsed)
}

func handleGet(w http.ResponseWriter, req *http.Request) error {
	u, _ := url.Parse(req.RequestURI)
	log.Println("GET requested ", u)

	pagectx := PageCtx{
		RootUrl: conf.Current.RootURLPattern,
		Buildnr: idl.Buildnr,
	}
	templName := "templates/vue/index.html"

	tmplIndex := template.Must(template.New("AppIndex").ParseFiles(templName))

	err := tmplIndex.ExecuteTemplate(w, "base", pagectx)
	if err != nil {
		return err
	}

	return nil
}

func handlePost(w http.ResponseWriter, req *http.Request) error {
	var err error
	lastPath := getURLForRoute(req.RequestURI)
	log.Println("Check the last path ", lastPath)
	switch lastPath {
	case "PlayURI":
		err = handlePlay(w, req)
	case "Pause":
		err = handlePause(w, req)
	case "ChangeVolume":
		err = handleChangeVolume(w, req)
	default:
		return fmt.Errorf("%s method is not supported", lastPath)
	}

	return err
}

func writeResponse(w http.ResponseWriter, resp interface{}) error {
	blobresp, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Write(blobresp)
	return nil
}

func writeErrorResponse(w http.ResponseWriter, errorcode int, resp interface{}) error {
	blobresp, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	http.Error(w, string(blobresp), errorcode)
	return nil
}
