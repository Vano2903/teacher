package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Host     string `yaml: "host"`
	Port     string `yaml: "port"`
	User     string `yaml: "user"`
	Password string `yaml: "password"`
	Dbname   string `yaml: "dbname"`
	Secret   string `yaml: "secret"`
}

var conf Config

func init() {
	//read the config.yaml, parse it and load the config struct
	host := os.Getenv("host")
	port := os.Getenv("port")
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")
	secret := os.Getenv("secret")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" || secret == "" {
		dat, err := ioutil.ReadFile("config_local.yaml")
		err = yaml.Unmarshal([]byte(dat), &conf)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	} else {
		conf.Dbname = dbname
		conf.Host = host
		conf.Port = port
		conf.User = user
		conf.Password = password
		conf.Secret = secret
	}

	fmt.Println(QueryTeacherByRegistrationNumber(2))
}
