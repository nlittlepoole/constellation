package main

import (
	"github.com/shibukawa/configdir"
	"path/filepath"
)

var osConfig configdir.ConfigDir

const APP_VENDOR = "Constellation"
const APP_NAME = "Observatory"

func init() {
	osConfig = configdir.New(APP_VENDOR, APP_NAME)
	// local path has highest priority
	osConfig.LocalPath, _ = filepath.Abs(".")
	if err := osConfig.QueryCacheFolder().MkdirAll(); err != nil {
		panic(err)
	}
	if err := osConfig.QueryFolders(configdir.Global)[0].MkdirAll(); err != nil {
		panic(err)
	}
	log.Info("Settings stored at ", osConfig.QueryFolders(configdir.Global)[0].Path)
}

func getCachePath() string {
	return osConfig.QueryCacheFolder().Path
}

func getConfigPath() string {
	return osConfig.QueryFolders(configdir.Global)[0].Path
}
