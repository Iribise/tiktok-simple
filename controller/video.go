package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"tiktok_test/common"
	"tiktok_test/constants"
	"tiktok_test/db"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Basic     BasicResponse
	NextTime  int64          `json:"next_time"`
	VideoList []common.Video `json:"video_list"`
}

type PublishResponse BasicResponse

type ListResponse struct {
	Basic     BasicResponse
	VideoList []common.Video `json:"video_list"`
}

func Feed(c *gin.Context) {
	token := c.Query("token")
	var key string
	_, login := userLoginInfo[token]
	if login {
		key = fmt.Sprintf("userfav:%d", userLoginInfo[token].Id)
	}
	var videoinfos []db.Video
	db.DB.Model(&db.Video{}).Order("created_at desc").Limit(30).Find(&videoinfos)
	var video_list []common.Video
	ctx := context.Background()
	for i := 0; i < len(videoinfos); i++ {
		videoinfos[i].Video.Id = int32(videoinfos[i].ID)
		if login {
			if exist, _ := db.RDB.SIsMember(ctx, key, videoinfos[i].Video.Id).Result(); exist {
				videoinfos[i].Video.IsFav = true
			}
		}
		db.DB.Model(&db.UserInfo{}).Where("id = ?", videoinfos[i].AuthorId).Find(&videoinfos[i].Video.Author)
		video_list = append(video_list, videoinfos[i].Video)
	}
	c.JSON(http.StatusOK, FeedResponse{
		Basic:     BasicResponse{StatusCode: 0, StatusMsg: "success"},
		NextTime:  time.Now().Unix(),
		VideoList: video_list,
	})
}

func Publish(c *gin.Context) {
	token := c.PostForm("token")
	if _, exist := userLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, PublishResponse{
			StatusCode: -1, StatusMsg: "user didn't log in",
		})
		return
	}
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, PublishResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		})
		return
	}
	title := c.PostForm("title")
	log.Println(file.Filename)
	name := fmt.Sprintf("%s_%s", title, filepath.Base(file.Filename))
	path := filepath.Join("./public/", name)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.JSON(http.StatusOK, PublishResponse{
			StatusCode: -1, StatusMsg: err.Error(),
		})
		return
	}
	userinfo := userLoginInfo[token]
	videoinfo := common.Video{
		PlayURL:  fmt.Sprintf("%s/static/%s", constants.VideoURLPrefix, name),
		CoverURL: "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavCount: 0,
		ComCount: 0,
		IsFav:    false,
		Title:    title,
	}
	db.DB.Create(&db.Video{AuthorId: userinfo.Id, Video: videoinfo})
	c.JSON(http.StatusOK, PublishResponse{
		StatusCode: 0, StatusMsg: "upload success",
	})
}

func ListPublish(c *gin.Context) {
	token := c.Query("token")
	if _, exist := userLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, PublishResponse{
			StatusCode: -1, StatusMsg: "user didn't log in",
		})
		return
	}
	userinfo := userLoginInfo[token]
	var list []db.Video
	db.DB.Model(&db.Video{}).Where("author_id = ?", userinfo.Id).Find(&list)
	var video_list []common.Video
	for i := 0; i < len(list); i++ {
		list[i].Video.Id = int32(list[i].ID)
		list[i].Video.Author = userinfo
		video_list = append(video_list, list[i].Video)
	}
	c.JSON(http.StatusOK, ListResponse{
		Basic:     BasicResponse{StatusCode: 0, StatusMsg: "success"},
		VideoList: video_list,
	})
}
