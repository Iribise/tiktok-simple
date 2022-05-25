package main

import (
	"tiktok_test/controller"
	"tiktok_test/db"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()
	controller.Init()
	r := gin.Default()
	initRouter(r)
	r.Run()
}
