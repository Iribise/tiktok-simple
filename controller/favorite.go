package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"tiktok_test/common"
	"tiktok_test/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FavResponse BasicResponse

type FavListResponse struct {
	Basic     BasicResponse
	VideoList []common.Video `json:"video_list"`
}

type ComResponse struct {
	Basic   BasicResponse
	Comment common.Comment `json:"comment"`
}

type ComListResponse struct {
	Basic       BasicResponse
	CommentList []common.Comment `json:"comment_list"`
}

func Favorite(c *gin.Context) {
	//user_id := c.Query("user_id")
	token := c.Query("token")
	video_id, _ := strconv.Atoi(c.Query("video_id"))
	action_type := c.Query("action_type")
	if _, exist := userLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, FavResponse{StatusCode: -1, StatusMsg: "user didn't log in"})
		return
	}
	user_id := userLoginInfo[token].Id
	ctx := context.Background()
	if action, _ := strconv.Atoi(action_type); action == 1 {
		db.RDB.SAdd(ctx, fmt.Sprintf("userfav:%d", user_id), video_id)
		db.DB.Model(&db.Video{}).Where("id = ?", video_id).Update("fav_count", gorm.Expr("fav_count + ?", 1))
	} else {
		db.RDB.SRem(ctx, fmt.Sprintf("userfav:%d", user_id), video_id)
		db.DB.Model(&db.Video{}).Where("id = ?", video_id).Update("fav_count", gorm.Expr("fav_count - ?", 1))
	}
	c.JSON(http.StatusOK, FavResponse{StatusCode: 0, StatusMsg: "success"})
}

func FavList(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	if _, exist := userLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, FavResponse{StatusCode: -1, StatusMsg: "user didn't log in"})
		return
	}
	ctx := context.Background()
	vals := db.RDB.SMembers(ctx, fmt.Sprintf("userfav:%s", user_id)).Val()
	if len(vals) == 0 {
		c.JSON(http.StatusOK, FavListResponse{
			Basic:     BasicResponse{StatusCode: 0, StatusMsg: "success"},
			VideoList: []common.Video{},
		})
	}
	videoids := make([]int, len(vals))
	for i := 0; i < len(vals); i++ {
		videoids[i], _ = strconv.Atoi(vals[i])
	}
	video_list := make([]common.Video, len(vals))
	var videos []db.Video
	db.DB.Model(&db.Video{}).Where(videoids).Find(&videos)
	for i := 0; i < len(videos); i++ {
		videos[i].Video.Id = int32(videos[i].ID)
		videos[i].Video.IsFav = true
		db.DB.Model(&db.UserInfo{}).Where("id = ?", videos[i].AuthorId).Find(&videos[i].Video.Author)
		video_list[i] = videos[i].Video
	}
	c.JSON(http.StatusOK, FavListResponse{
		Basic:     BasicResponse{StatusCode: 0, StatusMsg: "success"},
		VideoList: video_list,
	})
}

func Comment(c *gin.Context) {
	token := c.Query("token")
	if _, exist := userLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, FavResponse{StatusCode: -1, StatusMsg: "user didn't log in"})
		return
	}
	video_id, _ := strconv.Atoi(c.Query("video_id"))
	video_key := fmt.Sprintf("videocom:%d", video_id)
	user_key := fmt.Sprintf("usercom:%d", userLoginInfo[token].Id)
	ctx := context.Background()
	var comment_db db.Comment
	if action, _ := strconv.Atoi(c.Query("action_type")); action == 1 {
		comment_text := c.Query("comment_text")
		comment_db = db.Comment{
			Comment: common.Comment{
				Content: comment_text,
			},
			UserId: userLoginInfo[token].Id,
		}
		db.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&db.Comment{}).Create(&comment_db).Error; err != nil {
				return err
			}
			if err := tx.Model(&db.Comment{}).Order("created_at desc").Find(&comment_db).Error; err != nil {
				return err
			}
			return nil
		})
		db.DB.Model(&db.Video{}).Where("id = ?", video_id).Update("com_count", gorm.Expr("com_count + ?", 1))
		db.RDB.SAdd(ctx, video_key, int(comment_db.ID))
		db.RDB.SAdd(ctx, user_key, int(comment_db.ID))
	} else {
		comment_id, _ := strconv.Atoi(c.Query("comment_id"))
		db.DB.Model(&db.Comment{}).Where("id = ?", comment_id).Find(&comment_db)
		db.DB.Delete(&db.Comment{}, comment_id)
		db.DB.Model(&db.Video{}).Where("id = ?", video_id).Update("com_count", gorm.Expr("com_count - ?", 1))
		db.RDB.SRem(ctx, video_key, comment_id)
		db.RDB.SRem(ctx, user_key, comment_id)
	}
	comment_db.Comment.Id = int32(comment_db.ID)
	db.DB.Model(&db.UserInfo{}).Where("id = ?", comment_db.UserId).Find(&comment_db.Comment.User)
	comment_db.Comment.CreateDate = comment_db.CreatedAt.Format("01-02")
	c.JSON(http.StatusOK, ComResponse{
		Basic:   BasicResponse{StatusCode: 0, StatusMsg: "success"},
		Comment: comment_db.Comment,
	})
}

func ComList(c *gin.Context) {
	token := c.Query("token")
	if _, exist := userLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, FavResponse{StatusCode: -1, StatusMsg: "user didn't log in"})
		return
	}
	video_id, _ := strconv.Atoi(c.Query("video_id"))
	ctx := context.Background()
	vals := db.RDB.SMembers(ctx, fmt.Sprintf("videocom:%d", video_id)).Val()
	if len(vals) == 0 {
		c.JSON(http.StatusOK, ComListResponse{
			Basic:       BasicResponse{StatusCode: 0, StatusMsg: "success"},
			CommentList: []common.Comment{},
		})
		return
	}
	var comids []int = make([]int, len(vals))
	for i := 0; i < len(vals); i++ {
		comids[i], _ = strconv.Atoi(vals[i])
	}
	var comments []db.Comment
	comment_list := make([]common.Comment, len(comids))
	db.DB.Model(&db.Comment{}).Where(comids).Order("created_at desc").Find(&comments)
	for i := 0; i < len(comments); i++ {
		comments[i].Comment.Id = int32(comments[i].ID)
		comments[i].Comment.CreateDate = comments[i].CreatedAt.Format("01-02")
		db.DB.Model(&db.UserInfo{}).Where("id = ?", comments[i].UserId).Find(&comments[i].Comment.User)
		comment_list[i] = comments[i].Comment
	}
	c.JSON(http.StatusOK, ComListResponse{
		Basic:       BasicResponse{StatusCode: 0, StatusMsg: "success"},
		CommentList: comment_list,
	})
}
