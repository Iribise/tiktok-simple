package main

import (
	"github.com/Iribise/tiktok-simple/controller"
	"github.com/Iribise/tiktok-simple/db"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()
	controller.Init()
	r := gin.Default()
	initRouter(r)
	r.Run()
}
