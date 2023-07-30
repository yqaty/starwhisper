package mydb

import (
	//"cmd/internal/mydb/tools"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/spf13/viper"
)

type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uint      `json:"user_id" gorm:"index"`
	Username  string    `json:"username"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

type PostWithTags struct {
	Post
	Tags []string `json:"tags"`
}

func AddPost(pt *Post) error {
	log.Println("AddPost")
	db := GetDB()
	result := db.Create(&pt)
	return result.Error
}

func GetPostByID(post_id uint) (Post, error) {
	log.Println("GetPostByID")
	var ps Post
	ps.ID = post_id
	db := GetDB()
	result := db.First(&ps)
	if result.RowsAffected == 0 {
		return Post{}, errors.New("the post does not exist")
	}
	return ps, nil
}

func GetPostWithTagByID(post_id uint) (PostWithTags, error) {
	log.Println("GetPostWithTagByID")
	post, err := GetPostByID(post_id)
	if err != nil {
		return PostWithTags{}, err
	}
	tags, err := GetPostTags(post_id)
	if err != nil {
		return PostWithTags{}, err
	}
	return PostWithTags{post, tags}, nil
}

func GetUserPostNumber(user_id uint) (int64, error) {
	db := GetDB()
	var count int64
	db.Model(&Post{}).Where("user_id = ?", user_id).Count(&count)
	return count, nil
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

// descending order
func GetUserPosts(usrid uint, fr int, to int) ([]Post, error) {
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

func GetRandPost() (Post, error) {
	var post Post
	var count int64
	db := GetDB()
	result := db.Model(&Post{}).Count(&count)
	if result.Error != nil {
		return Post{}, result.Error
	}
	rand.Seed(time.Now().Unix())
	index := rand.Int63n(count)
	result = db.Offset(int(index)).First(&post)
	if result.Error != nil {
		return Post{}, result.Error
	}
	return post, nil
}

func Posts2PostsWithTags(posts []Post) ([]PostWithTags, error) {
	postswithtags := make([]PostWithTags, len(posts))
	for i, post := range posts {
		tags, _ := GetPostTags(post.ID)
		postswithtags[i] = PostWithTags{post, tags}
	}
	return postswithtags, nil
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
}
