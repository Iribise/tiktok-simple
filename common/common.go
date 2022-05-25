package common

type UserInfo struct {
	Id            int32  `gorm:"primaryKey" json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int32  `json:"follow_count,omitempty"`
	FollowerCount int32  `json:"follower_count,omitempty"`
	IsFollow      bool   `gorm:"-" json:"is_follow,omitempty"`
}

type User struct {
	Id       int32 `gorm:"primaryKey"`
	Password string
}

type Video struct {
	Id       int32    `gorm:"-" json:"id"`
	Author   UserInfo `gorm:"-" json:"author"`
	PlayURL  string   `json:"play_url,omitempty"`
	CoverURL string   `json:"cover_url,omitempty"`
	FavCount int32    `json:"favorite_count,omitempty"`
	ComCount int32    `json:"comment_count,omitempty"`
	IsFav    bool     `gorm:"-" json:"is_favorite,omitempty"`
	Title    string   `json:"title,omitempty"`
}

type Comment struct {
	Id         int32    `gorm:"-" json:"id"`
	User       UserInfo `gorm:"-" json:"user"`
	Content    string   `json:"content"`
	CreateDate string   `gorm:"-" json:"create_date"`
}
