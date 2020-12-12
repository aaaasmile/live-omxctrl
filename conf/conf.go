package conf

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ServiceURL      string
	RootURLPattern  string
	UseRelativeRoot bool
	DebugVerbose    bool
	OmxCmdParams    string
	DBPath          string
	TmpInfo         string
	VueLibName      string
	SoundCloud      *SoundCloud
}

type SoundCloud struct {
	ClientID  string
	AuthToken string
	UserAgent string
}

var Current = &Config{}

func ReadConfig(configfile string) *Config {
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := toml.DecodeFile(configfile, &Current); err != nil {
		log.Fatal(err)
	}
	// TODO read SoundCloud from soundclod.json plugins dir.
	return Current
}
