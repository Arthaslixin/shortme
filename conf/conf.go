package conf

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type http struct {
	Listen string `toml:"listen"`
}

type common struct {
	BlackShortURLs    []string `toml:"black_short_urls"`
	BlackShortURLsMap map[string]bool
	BaseString        string `toml:"base_string"`
	BaseStringLength  uint64
	DomainName        string `toml:"domain_name"`
	Schema            string `toml:"schema"`
}

type config struct {
	Http   http   `toml:"http"`
	Common common `toml:"common"`
}

var Conf config

func MustParseConfig(configFile string) {
	if fileInfo, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			log.Panicf("configuration file %v does not exist.", configFile)
		} else {
			log.Panicf("configuration file %v can not be stated. %v", configFile, err)
		}
	} else {
		if fileInfo.IsDir() {
			log.Panicf("%v is a directory name", configFile)
		}
	}

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Panicf("read configuration file error. %v", err)
	}
	content = bytes.TrimSpace(content)

	err = toml.Unmarshal(content, &Conf)
	if err != nil {
		log.Panicf("unmarshal toml object error. %v", err)
	}

	// short url black list
	Conf.Common.BlackShortURLsMap = make(map[string]bool)
	for _, blackShortURL := range Conf.Common.BlackShortURLs {
		Conf.Common.BlackShortURLsMap[blackShortURL] = true
	}

	// base string
	Conf.Common.BaseStringLength = uint64(len(Conf.Common.BaseString))
}
