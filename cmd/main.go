package main

import (
	"log"

	"github.com/dunkykorZhik/avito-tech/config"
	"github.com/dunkykorZhik/avito-tech/internal/db"
	"github.com/dunkykorZhik/avito-tech/internal/server"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("Cannot get the Config", err)
		return
	}

	db, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Cannot connect to Db ", err)
		return
	}
	if db == nil {
		log.Fatal("Cannot connect to Db ", err)
		return
	}
	log.Println("Got Db")

	server.NewFiberServer(cfg, db).Start()

}
