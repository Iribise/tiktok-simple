package db

import (
	"github.com/Iribise/tiktok-simple/common"
	"github.com/Iribise/tiktok-simple/constants"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RDB *redis.Client

func Init() {
	var err error
	DB, err = gorm.Open(mysql.Open(constants.MysqlDSN), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err.Error())
	}
	DB.AutoMigrate(&UserInfo{}) // skip err check for simplicity
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Video{})
	DB.AutoMigrate(&Comment{})

	RDB = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       0,
	})
}

type UserInfo common.UserInfo

type User common.User

type Video struct {
	gorm.Model
	common.Video
	AuthorId int32
}

type Comment struct {
	gorm.Model
	common.Comment
	UserId int32
}
