package main

type Redisee struct {
	Config Config
}

type Config struct {
	Host        string
	Port        int
	Password    string
	Db          int
	Separator   string
	ScanPattern string
	Concurrency int
}

func (*Redisee) Run() {

}
