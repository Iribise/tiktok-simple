package main

import (
	"tiktok_test/controller"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	r.Static("/static", "./public")
	apiRouter := r.Group("/douyin")
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.GET("/user/", controller.GetUserInfo)
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.ListPublish)
	apiRouter.POST("/favorite/action/", controller.Favorite)
	apiRouter.GET("/favorite/list/", controller.FavList)
	apiRouter.POST("/comment/action/", controller.Comment)
	apiRouter.GET("/comment/list/", controller.ComList)
	apiRouter.POST("/relation/action/", controller.Follow)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
}
