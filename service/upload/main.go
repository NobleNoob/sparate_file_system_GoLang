package main

import (
	cfg "filestore-server/config"
	"filestore-server/route"
)

func main() {

	route:=route.Router()
	route.Run(cfg.UploadServiceHost)
}