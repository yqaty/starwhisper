package mydb

import (
	//"cmd/internal/mydb/tools"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserID        uint   `json:"-"`
	UserName      string `json:"user_name"`
	CommentNumber uint   `json:"comment_number"`
	Title         string `json:"title"`
	Context       string `json:"context"`
	Link          string `json:"link"`
}

var _domain string

func AddPost(pt *Post) error {
	log.Println("AddPost")
	db := GetDB()
	result := db.Create(&pt)
	if result.Error != nil {
		return result.Error
	}
	var err error
	pt.UserName, err = GetAnonymity(pt.UserID, pt.ID)
	if err != nil {
		return err
	}
	pt.Link = _domain + "/article/" + strconv.Itoa(int(pt.ID))
	result = db.Save(&pt)
	return result.Error
}

func UpdatePost(pt *Post) error {
	log.Println("UpdatePost")
	var dbpt Post
	db := GetDB()
	result := db.First(&dbpt, pt.ID)
	if result.Error != nil {
		return result.Error
	}
	if dbpt.UserID != pt.UserID {
		return errors.New("permission denied")
	}
	dbpt.Title = pt.Title
	dbpt.Context = pt.Context
	result = db.Save(&dbpt)
	*pt = dbpt
	return result.Error
}

func GetPost(postid uint) (Post, error) {
	log.Println("GetPost")
	var ps Post
	ps.ID = postid
	db := GetDB()
	result := db.First(&ps)
	if result.RowsAffected == 0 {
		return Post{}, errors.New("the post does not exist")
	}
	return ps, nil
}

func SearchPosts(str string, fr time.Time, to time.Time, st int, ed int) ([]Post, error) {
	log.Println("SearchPost")
	log.Println(str, fr, to, st, ed)
	var ps []Post
	db := GetDB()
	result := db.Where("title LIKE ? AND created_at BETWEEN ? AND ?", "%"+str+"%", fr, to).Limit(ed - st + 1).Offset(st - 1).Find(&ps)
	if result.Error != nil {
		return nil, result.Error
	}
	return ps, nil
}

//descending order
func UserPosts(usrid uint, fr int, to int) ([]Post, error) {
	log.Println("UserPosts")
	var pts []Post
	db := GetDB()
	log.Println(usrid, fr, to)
	result := db.Where("user_id = ?", usrid).Order("id desc").Limit(to - fr + 1).Offset(fr - 1).Find(&pts)
	if result.Error != nil {
		return nil, result.Error
	}
	return pts, nil
}

func DelPost(pt *Post) error {
	log.Println("DelPost")
	var dbpt Post
	db := GetDB()
	result := db.First(&dbpt, pt.ID)
	if result.Error != nil {
		return result.Error
	}
	if pt.UserID != dbpt.UserID {
		return errors.New("permission denied")
	}
	result = db.Delete(&Post{}, pt.ID)
	return result.Error
}

func init() {
	db := GetDB()
	db.AutoMigrate(&Post{})
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("%w", err))
	}
	_domain = viper.GetString("domain") + ":" + viper.GetString("port")
}
