package config

const ALL_DB = -1

var Config = ConfigT{}

type ConfigT struct {
	Host        string
	Port        int
	Password    string
	Db          int
	Separator   string
	ScanPattern string
	Concurrency int
}
