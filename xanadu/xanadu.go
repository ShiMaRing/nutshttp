package main

import (
	"github.com/xujiajun/nutsdb"
	"log"
	"nutshttp"
)

func main() {
	opt := nutsdb.DefaultOptions
	opt.Dir = "./data"

	db, err := nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		db.Close()
	}()

	go func() {
		if err := nutshttp.Enable(db); err != nil {
			panic(err)
		}
	}()

	select {}
}
