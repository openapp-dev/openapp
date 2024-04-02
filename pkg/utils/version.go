package utils

var (
	gitVersion = "latest"
	gitCommit  = "unknown"
)

type OpenAPPVersion struct {
	GitVersion string `json:"gitVersion"`
	GitCommit  string `json:"gitCommit"`
}

func GetOpenAPPVersion() *OpenAPPVersion {
	v := &OpenAPPVersion{
		GitVersion: gitVersion,
		GitCommit:  gitCommit,
	}
	return v
}
