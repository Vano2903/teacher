package main

import (
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

var (
	corsiTable string = `
	CREATE TABLE IF NOT EXISTS corsi (
  		ID int NOT NULL AUTO_INCREMENT,
		materia varchar(45) NOT NULL,
		api_value int NOT NULL,
		PRIMARY KEY (ID)
	) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`
	esamiTable string = `
	CREATE TABLE IF NOT EXISTS esami (
		ID int NOT NULL AUTO_INCREMENT,
		difficolta varchar(10) NOT NULL,
		ID_corso int DEFAULT NULL,
		numero_domande int NOT NULL,
		nome varchar(45) NOT NULL,
		ID_insegnante int NOT NULL,
		PRIMARY KEY (ID)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`
	insegnaTable string = `
	CREATE TABLE IF NOT EXISTS insegna (
		ID int NOT NULL,
		id_insegnante int DEFAULT NULL,
		id_corso int DEFAULT NULL,
		PRIMARY KEY (ID)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`
	insegnanteTable string = `
	CREATE TABLE IF NOT EXISTS insegnante (
		ID int NOT NULL AUTO_INCREMENT,
		nome varchar(45) NOT NULL,
		cognome varchar(45) NOT NULL,
		matricola int NOT NULL,
		password char(64) NOT NULL,
		PRIMARY KEY (ID)
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`
	risultatiEsamiTable string = `
	CREATE TABLE IF NOT EXISTS risultati_esami (
		ID int NOT NULL AUTO_INCREMENT,
		nome_studente varchar(45) NOT NULL,
		cognome_studente varchar(45) NOT NULL,
		contenuto mediumtext NOT NULL,
		ID_esame int NOT NULL,
		tentativi int NOT NULL,
		PRIMARY KEY (ID)
	)ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`
)

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

	db, err := ConnectToDb()
	if err != nil {
		log.Fatalf("error connecting to the bd: " + err.Error())
	}

	log.Println("connection with db established")
	defer db.Close()

	_, err = db.Exec(insegnanteTable)
	if err != nil {
		log.Fatalf("insegnante table creation failed: %s", err.Error())
	}

	_, err = db.Exec(insegnaTable)
	if err != nil {
		log.Fatalf("insegna table creation failed: %s", err.Error())
	}

	_, err = db.Exec(esamiTable)
	if err != nil {
		log.Fatalf("esami table creation failed: %s", err.Error())
	}

	_, err = db.Exec(corsiTable)
	if err != nil {
		log.Fatalf("corsi table creation failed: %s", err.Error())
	}

	_, err = db.Exec(risultatiEsamiTable)
	if err != nil {
		log.Fatalf("risultati-esami table creation failed: %s", err.Error())
	}

}
