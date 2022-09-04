package main

import (
	"log"
	db_saver_service "workspace/dbsavergo/kitex_gen/db_saver_service/dbsaverservice"
)

func main() {
	svr := db_saver_service.NewServer(new(DBSaverServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
