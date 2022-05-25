package controller

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/Iribise/tiktok-simple/common"
	"github.com/Iribise/tiktok-simple/db"
	"github.com/gin-gonic/gin"
)

var userLoginInfo map[string]common.UserInfo

var nextId int32

func Init() {
	userLoginInfo = make(map[string]common.UserInfo)
	db.DB.Raw("SELECT MAX(id) FROM users").Scan(&nextId)
	nextId++
}

type RegisterResponse struct {
	Basic BasicResponse
	Id    int32  `json:"user_id,omitempty"`
	Token string `json:"token,omitempty"`
}

type LoginResponse RegisterResponse

type UserInfoResponse struct {
	Basic BasicResponse
	Info  common.UserInfo
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	hash, err := hashPassword(password)
	if err != nil {
		c.JSON(http.StatusOK, RegisterResponse{
			Basic: BasicResponse{StatusCode: -1, StatusMsg: "register failed"},
		})
		return
	}
	token := username + hash
	if _, exist := userLoginInfo[token]; exist {
		c.JSON(http.StatusOK, RegisterResponse{
			Basic: BasicResponse{StatusCode: -1, StatusMsg: "user already exist"},
		})
		return
	}
	var userinfo db.UserInfo
	result := db.DB.Where("name = ?", username).Take(&userinfo)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, RegisterResponse{
			Basic: BasicResponse{StatusCode: -1, StatusMsg: "user already exist"},
		})
		return
	}
	userinfo = db.UserInfo{Id: nextId, Name: username, FollowCount: 0, FollowerCount: 0, IsFollow: false}
	db.DB.Create(&userinfo)
	user := db.User{Id: nextId, Password: hash}
	db.DB.Create(&user)
	userLoginInfo[token] = common.UserInfo(userinfo)
	c.JSON(http.StatusOK, RegisterResponse{
		Basic: BasicResponse{StatusCode: 0},
		Id:    nextId,
		Token: token,
	})
	nextId++
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	var userinfo db.UserInfo
	result := db.DB.Where("name = ?", username).Take(&userinfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, LoginResponse{
			Basic: BasicResponse{StatusCode: -1, StatusMsg: "user doesn't exist"},
		})
		return
	}
	var user db.User
	db.DB.Where("id = ?", userinfo.Id).Take(&user)
	if !checkPassword(password, user.Password) {
		c.JSON(http.StatusOK, LoginResponse{
			Basic: BasicResponse{StatusCode: -1, StatusMsg: "wrong password"},
		})
		return
	}
	token := username + user.Password
	if _, exist := userLoginInfo[token]; !exist {
		userLoginInfo[token] = common.UserInfo(userinfo)
	}
	c.JSON(http.StatusOK, LoginResponse{
		Basic: BasicResponse{StatusCode: 0},
		Id:    user.Id,
		Token: token,
	})
}

func GetUserInfo(c *gin.Context) {
	token := c.Query("token")
	if userinfo, exist := userLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserInfoResponse{
			Basic: BasicResponse{StatusCode: 0},
			Info:  common.UserInfo(userinfo),
		})
	} else {
		c.JSON(http.StatusOK, UserInfoResponse{Basic: BasicResponse{StatusCode: -1, StatusMsg: "user didn't log in"}})
	}
}
