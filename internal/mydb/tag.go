package mydb

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

type Tag struct {
	ID  uint   `json:"id" gorm:"primaryKey:autoIncrement"`
	Tag string `json:"tag"  gorm:"index"`
}

type PostTag struct {
	PostID uint `gorm:"primaryKey"`
	TagID  uint `gorm:"primaryKey"`
}

func AddTag(tag string) error {
	log.Println("AddTag")
	db := GetDB()
	result := db.Save(&Tag{Tag: tag})
	return result.Error
}

func Tag2IDs(tags []string) ([]uint, error) {
	log.Println("Tag2IDs")
	tag_ids := make([]uint, len(tags))
	var tag_id Tag
	db := GetDB()
	for i, tag := range tags {
		tag_id = Tag{0, tag}
		result := db.Where("tag = ?", tag).First(&tag_id)
		if result.Error != nil {
			return nil, result.Error
		}
		tag_ids[i] = tag_id.ID
	}
	return tag_ids, nil
}

func ID2Tags(tag_ids []uint) ([]string, error) {
	log.Println("ID2Tags")
	tags := make([]string, len(tag_ids))
	var tag Tag
	db := GetDB()
	for i, id := range tag_ids {
		tag = Tag{id, ""}
		result := db.First(&tag)
		if result.Error != nil {
			return nil, result.Error
		}
		tags[i] = tag.Tag
	}
	return tags, nil
}

func AddTags(post_id uint, tags []string) error {
	log.Println("AddTags")
	tag_ids, err := Tag2IDs(tags)
	if err != nil {
		return nil
	}
	db := GetDB()
	tx := db.Begin()
	for _, id := range tag_ids {
		result := tx.Create(&PostTag{post_id, id})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}
	return tx.Commit().Error
}

func GetPostByTags(tags []string) (uint, error) {
	log.Println("GetPostByTags", tags)
	db := GetDB()
	tag_ids, err := Tag2IDs(tags)
	if err != nil {
		return 0, nil
	}
	var posts []PostTag
	log.Println(tag_ids)
	db.Model(&PostTag{}).
		Select("post_id").
		Where("tag_id IN (?)", tag_ids).
		Group("post_id").
		Having("COUNT(DISTINCT tag_id) = ?", len(tags)).
		Find(&posts)
	if len(posts) == 0 {
		return 0, errors.New("post is not exist")
	}
	rand.Seed(time.Now().Unix())
	index := rand.Intn(len(posts))
	return posts[index].PostID, nil
}

func DelPostTags(post_id uint) error {
	log.Println("DelPostTags")
	db := GetDB()
	db.Where("post_id = ?", post_id).Delete(&PostTag{})
	return nil
}

func GetPostTags(post_id uint) ([]string, error) {
	log.Println("GetPostTags")
	var tags []PostTag
	db := GetDB()
	db.Where("post_id = ?", post_id).Find(&tags)
	var tag_ids = make([]uint, len(tags))
	for i, tag := range tags {
		tag_ids[i] = tag.TagID
	}
	return ID2Tags(tag_ids)
}

func init() {
	db := GetDB()
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&PostTag{})
	db.Save(&Tag{1, "记录思考"})
	db.Save(&Tag{2, "倾诉烦恼"})
	db.Save(&Tag{3, "分享趣事"})
	db.Save(&Tag{4, "找寻另一个自己"})
	db.Save(&Tag{5, "悄悄话"})
	db.Save(&Tag{6, "感情里的那些事"})
	db.Save(&Tag{7, "祈福许愿"})
	db.Save(&Tag{8, "无聊ing"})
}
