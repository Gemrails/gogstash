package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"regexp"

	"github.com/codegangsta/inject"
	"github.com/tsaikd/KDGoLib/errutil"
)

type Config struct {
	inject.Injector `json:"-"`
	InputRaw        []ConfigRaw `json:"input,omitempty"`
	FilterRaw       []ConfigRaw `json:"filter,omitempty"`
	OutputRaw       []ConfigRaw `json:"output,omitempty"`
}

func LoadFromFile(path string) (config Config, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		err = errutil.New("Failed to read config file, "+path, err)
		return
	}

	return LoadFromData(data)
}

func LoadFromString(text string) (config Config, err error) {
	return LoadFromData([]byte(text))
}

func LoadFromData(data []byte) (config Config, err error) {
	if data, err = StripComments(data); err != nil {
		err = errutil.New("Failed to strip comments from json", err)
		return
	}

	if err = json.Unmarshal(data, &config); err != nil {
		err = errutil.New("Failed unmarshalling json", err)
		return
	}

	config.Injector = inject.New()
	return
}

func ReflectConfig(confraw *ConfigRaw, conf interface{}) (err error) {
	data, err := json.Marshal(confraw)
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, conf); err != nil {
		return
	}

	return
}

// Supported comment formats
// format 1: ^\s*#
// format 2: ^\s*//
func StripComments(data []byte) (out []byte, err error) {
	reForm1 := regexp.MustCompile(`^\s*#`)
	reForm2 := regexp.MustCompile(`^\s*//`)
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0) // Windows
	lines := bytes.Split(data, []byte("\n"))
	filtered := make([][]byte, 0)

	for _, line := range lines {
		if reForm1.Match(line) {
			continue
		}
		if reForm2.Match(line) {
			continue
		}
		filtered = append(filtered, line)
	}

	out = bytes.Join(filtered, []byte("\n"))
	return
}
