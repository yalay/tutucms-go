package conf

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
)

var gConfig = &Config{}

type Config struct {
	Author     string
	UserId     int32
	ComeUrl    string
	DictPaths  []string
	DbUser     string
	DbPassword string
	DbName     string
}

func init() {
	var configFile string
	flag.StringVar(&configFile, "c", "conf/config.yaml", "conf file path")
	flag.Parse()

	reloadConf(configFile)
}

func GetAuthor() string {
	return gConfig.Author
}

func GetComeUrl() string {
	return gConfig.ComeUrl
}

func GetUserId() int32 {
	return gConfig.UserId
}

func GetDictPath() string {
	return strings.Join(gConfig.DictPaths, ",")
}

func GetDbUser() string {
	return gConfig.DbUser
}

func GetDbPassword() string {
	return gConfig.DbPassword
}

func GetDbName() string {
	return gConfig.DbName
}

func reloadConf(configFile string) {
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Panicf("read config file err:%v\n", err)
	}

	err = yaml.Unmarshal(configData, gConfig)
	if err != nil {
		log.Panicf("parse config file err:%v\n", err)
	}

	log.Printf("config:%+v", gConfig)
	go reloadYamlFile(configFile, time.Minute, gConfig)
}

func reloadYamlFile(configFile string, duration time.Duration, serverConf *Config) {
	var lastMtime = getFileMtime(configFile)
	for {
		time.Sleep(duration)
		if curMtime := getFileMtime(configFile); curMtime > lastMtime {
			lastMtime = curMtime
			configData, err := ioutil.ReadFile(configFile)
			if err != nil {
				log.Panicf("read config file err:%v\n", err)
			}
			err = yaml.Unmarshal(configData, &serverConf)
			if err != nil {
				log.Panicf("parse config file err:%v\n", err)
			}
			log.Printf("config:%+v", serverConf)
		}
	}
}

func getFileMtime(file string) int64 {
	fileInfo, err := os.Stat(file)
	if err != nil {
		log.Fatalf("file stat err:%v\n", err)
		return 0
	}
	return fileInfo.ModTime().Unix()
}
