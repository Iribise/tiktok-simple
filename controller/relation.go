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

type FollowResponse BasicResponse

type FollowListResponse struct {
	Basic    BasicResponse
	UserList []common.UserInfo `json:"user_list"`
}

type FollowerResponse FollowResponse

func Follow(c *gin.Context) {
	token := c.Query("token")
	if _, exist := userLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, FollowResponse{StatusCode: -1, StatusMsg: "user didn't log in"})
		return
	}
	user_id := userLoginInfo[token].Id
	to_user_id, _ := strconv.Atoi(c.Query("to_user_id"))
	user_key := fmt.Sprintf("userfollow:%d", user_id)
	to_user_key := fmt.Sprintf("userfollower:%d", to_user_id)
	ctx := context.Background()
	if action, _ := strconv.Atoi(c.Query("action_type")); action == 1 {
		db.RDB.SAdd(ctx, user_key, to_user_id)
		db.RDB.SAdd(ctx, to_user_key, user_id)
		db.DB.Model(&db.UserInfo{}).Where("id = ?", user_id).Update("folllow_count", gorm.Expr("folllow_count + ?", 1))
		db.DB.Model(&db.UserInfo{}).Where("id = ?", to_user_id).Update("folllower_count", gorm.Expr("folllower_count + ?", 1))
	} else {
		db.RDB.SRem(ctx, user_key, to_user_id)
		db.RDB.SRem(ctx, to_user_key, user_id)
		db.DB.Model(&db.UserInfo{}).Where("id = ?", user_id).Update("folllow_count", gorm.Expr("folllow_count - ?", 1))
		db.DB.Model(&db.UserInfo{}).Where("id = ?", to_user_id).Update("folllower_count", gorm.Expr("folllower_count - ?", 1))
	}
	c.JSON(http.StatusOK, FollowResponse{StatusCode: 0, StatusMsg: "success"})
}

func FollowList(c *gin.Context) {
	token := c.Query("token")
	if _, exist := userLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, FollowListResponse{
			Basic:    BasicResponse{StatusCode: -1, StatusMsg: "user didn't log in"},
			UserList: []common.UserInfo{},
		})
		return
	}
	user_id := userLoginInfo[token].Id
	ctx := context.Background()
	vals := db.RDB.SMembers(ctx, fmt.Sprintf("userfollow:%d", user_id)).Val()
	if len(vals) == 0 {
		c.JSON(http.StatusOK, FollowListResponse{
			Basic:    BasicResponse{StatusCode: 0, StatusMsg: "success"},
			UserList: []common.UserInfo{},
		})
	}
	userids := make([]int, len(vals))
	for i := 0; i < len(vals); i++ {
		userids[i], _ = strconv.Atoi(vals[i])
	}
	var follows []db.UserInfo
	follow_list := make([]common.UserInfo, len(vals))
	db.DB.Model(&db.UserInfo{}).Where(userids).Find(&follows)
	for i := 0; i < len(vals); i++ {
		follows[i].IsFollow = true
		follow_list[i] = common.UserInfo(follows[i])
	}
	c.JSON(http.StatusOK, FollowListResponse{
		Basic:    BasicResponse{StatusCode: 0, StatusMsg: "success"},
		UserList: follow_list,
	})
}

func FollowerList(c *gin.Context) {
	token := c.Query("token")
	if _, exist := userLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, FollowListResponse{
			Basic:    BasicResponse{StatusCode: -1, StatusMsg: "user didn't log in"},
			UserList: []common.UserInfo{},
		})
		return
	}
	user_id := userLoginInfo[token].Id
	ctx := context.Background()
	vals := db.RDB.SMembers(ctx, fmt.Sprintf("userfollower:%d", user_id)).Val()
	if len(vals) == 0 {
		c.JSON(http.StatusOK, FollowListResponse{
			Basic:    BasicResponse{StatusCode: 0, StatusMsg: "success"},
			UserList: []common.UserInfo{},
		})
	}
	userids := make([]int, len(vals))
	for i := 0; i < len(vals); i++ {
		userids[i], _ = strconv.Atoi(vals[i])
	}
	var followers []db.UserInfo
	follower_list := make([]common.UserInfo, len(vals))
	db.DB.Model(&db.UserInfo{}).Where(userids).Find(&followers)
	for i := 0; i < len(vals); i++ {
		followers[i].IsFollow = true
		follower_list[i] = common.UserInfo(followers[i])
	}
	c.JSON(http.StatusOK, FollowListResponse{
		Basic:    BasicResponse{StatusCode: 0, StatusMsg: "success"},
		UserList: follower_list,
	})
}
