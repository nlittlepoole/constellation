package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Settings struct {
	Driver        string  `yaml:"driver"`
	Location      string  `yaml:"location"`
	Threshold     int64   `yaml:"threshold"`
	SampleRate    float64 `yaml:"sample_rate"`
	SessionUnit   string  `yaml:"session_unit"`
	SessionLength int64   `yaml:"session_length"`
}

var DEFAULT_SETTINGS Settings = Settings{
	Driver:        "mon0",
	Location:      "Wayne Manor",
	Threshold:     100,
	SampleRate:    1.0,
	SessionUnit:   "minute",
	SessionLength: 30,
}

const SETTINGS_FILE string = "settings.yaml"

var ACTIVE_SETTINGS Settings

func settingsPath() string {
	return filepath.Join(getConfigPath(), SETTINGS_FILE)
}

func init() {
	if _, err := os.Stat(settingsPath()); os.IsNotExist(err) {
		ACTIVE_SETTINGS = DEFAULT_SETTINGS
		err := ACTIVE_SETTINGS.Save()
		if err != nil {
			panic(err)
		}
	} else {
		dat, _ := ioutil.ReadFile(settingsPath())
		err := yaml.Unmarshal(dat, &ACTIVE_SETTINGS)
		if err != nil {
			panic(err)
		}
	}
}

func (s *Settings) Session() time.Duration {
	var unit time.Duration
	switch s.SessionUnit {
	case "hour":
		unit = time.Hour
	case "minute":
		unit = time.Minute
	default:
		unit = time.Hour * 24
	}
	return unit * time.Duration(s.SessionLength)
}

func (s *Settings) Save() error {
	bytes, err := yaml.Marshal(s)
	if err == nil {
		err = ioutil.WriteFile(settingsPath(), bytes, 0644)
	}
	return err
}
