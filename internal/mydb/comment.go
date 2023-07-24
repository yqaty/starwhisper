package mydb

import (
	//"cmd/internal/mydb"
	//"cmd/internal/mydb/tools"
	"errors"
	"log"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID        uint   `json:"-"`
	UserName      string `json:"user_name"`
	PostID        uint   `json:"post_id"`
	BelongID      uint   `json:"belong_id"`
	CommentNumber uint   `json:"comment_number"`
	Floor         uint   `json:"floor"`
	Context       string `json:"context"`
}

//if insert a new comment, add one to the comment number of comments or posts it replied to, return its floor
func ChaCmtNum(cmt *Comment, db *gorm.DB) (uint, error) {
	log.Println("ChaCmtNum")
	if cmt.BelongID == 0 {
		var ps Post
		result := db.First(&ps, cmt.PostID)
		if result.Error != nil {
			return 0, result.Error
		}
		ps.CommentNumber++
		db.Save(&ps)
		return ps.CommentNumber, nil
	} else {
		var bel Comment
		result := db.First(&bel, cmt.BelongID)
		if result.Error != nil {
			return 0, result.Error
		}
		bel.CommentNumber++
		db.Save(&bel)
		return bel.CommentNumber, nil
	}
}

func AddComment(cmt *Comment) error {
	log.Println("AddComment")
	db := GetDB()
	tx := db.Begin()
	var err error
	cmt.Floor, err = ChaCmtNum(cmt, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	cmt.UserName, err = GetAnonymity(cmt.UserID, cmt.PostID)
	if err != nil {
		return err
	}
	result := tx.Create(&cmt)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	return tx.Commit().Error
}

func GetCommentID(cmtid uint) (Comment, error) {
	log.Println("GetComment")
	var cmt Comment
	db := GetDB()
	result := db.First(&cmt, cmtid)
	if result.RowsAffected == 0 {
		return Comment{}, errors.New("the comment does not exist")
	}
	return cmt, nil
}

//descending order
func UserComments(usrid uint, fr int, to int) ([]Comment, error) {
	log.Println("UserComments")
	var cmts []Comment
	db := GetDB()
	result := db.Where("user_id = ?", usrid).Order("id desc").Limit(to - fr + 1).Offset(fr - 1).Find(&cmts)
	if result.Error != nil {
		return nil, result.Error
	}
	return cmts, nil
}

func GetComment(post_id uint, belong_id uint, floor uint) (Comment, error) {
	log.Println("GetComment")
	var cmt Comment
	db := GetDB()
	result := db.Where("post_id = ? AND belong_id = ? AND floor = ?", post_id, belong_id, floor).First(&cmt)
	if result.Error != nil {
		return Comment{}, result.Error
	}
	return cmt, nil
}

func DelComment(cmt *Comment) error {
	log.Println("DelCommenrt")
	var dbcmt Comment
	db := GetDB()
	result := db.First(&dbcmt, cmt.ID)
	if result.Error != nil {
		return result.Error
	}
	if cmt.UserID != dbcmt.UserID {
		return errors.New("permission denied")
	}
	result = db.Delete(&Comment{}, cmt.ID)
	return result.Error
}

func init() {
	db := GetDB()
	db.AutoMigrate(&Comment{})
}
