package utils

import (
	"encoding/json"
)

var (
	gitVersion = "latest"
	gitCommit  = "unknown"
)

type OpenAPPVersion struct {
	GitVersion string `json:"gitVersion"`
	GitCommit  string `json:"gitCommit"`
}

func GetOpenAPPVersion() string {
	v := &OpenAPPVersion{
		GitVersion: gitVersion,
		GitCommit:  gitCommit,
	}
	data, _ := json.Marshal(v)
	return string(data)
}
