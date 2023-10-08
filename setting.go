package main

import "os"

type setting struct {
	ServerHost string
	ServerPort string
	PgHost     string
	PgPort     string
	PgUser     string
	PgPass     string
	PgBase     string
}

var cfg setting

func init() {
	file, err := os.Open("setting.cfg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

}
