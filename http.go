package nutshttp

import (
	"github.com/xujiajun/nutsdb"
	"log"
)

func Enable(db *nutsdb.DB) error {
	server, err := NewNutsHTTPServer(db)
	if err != nil {
		log.Fatalln(err)
	}

	return server.Run(":8299")
}
