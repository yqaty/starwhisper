package mydb

import (
	//"cmd/internal/mydb"
	//"cmd/internal/mydb/tools"
	"errors"
	"log"
	"time"
)

type Chat struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PostID    uint      `json:"post_id"`
	UserID1   uint      `json:"user_id1"`
	UserID2   uint      `json:"user_id2"`
	SendID    uint      `json:"send_id"`
	SendName  string    `json:"send_name"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

type Unseen struct {
	UserID uint `json:"user_id" gorm:"primaryKey"`
	PostID uint `json:"post_id" gorm:"primaryKey"`
	SendID uint `json:"send_id" gorm:"primaryKey"`
}

func AddChat(cht *Chat) error {
	log.Println("AddChat")
	log.Println(cht)
	db := GetDB()
	result := db.Create(cht)
	if result.Error != nil {
		return result.Error
	}
	db.Save(&Unseen{UserID: cht.UserID1 ^ cht.UserID2 ^ cht.SendID, PostID: cht.PostID, SendID: cht.SendID})
	return nil
}

func AddUnseen(us *Unseen) error {
	log.Println("AddUnseen")
	db := GetDB()
	db.Save(us)
	return nil
}

func DelUnseen(us *Unseen) error {
	log.Println("DelUnseen")
	db := GetDB()
	db.Delete(us)
	return nil
}

/*func GetCommentID(cmtid uint) (Comment, error) {
	log.Println("GetComment")
	var cmt Comment
	db := GetDB()
	result := db.First(&cmt, cmtid)
	if result.RowsAffected == 0 {
		return Comment{}, errors.New("the comment does not exist")
	}
	return cmt, nil
}

// descending order
func UserComments(usrid uint, fr int, to int) ([]Comment, error) {
	log.Println("UserComments")
	var cmts []Comment
	db := GetDB()
	result := db.Where("user_id = ?", usrid).Order("id desc").Limit(to - fr + 1).Offset(fr - 1).Find(&cmts)
	if result.Error != nil {
		return nil, result.Error
	}
	return cmts, nil
}*/

func GetPostChats(post_id uint, user_id uint, fr int, to int) ([]Chat, error) {
	log.Println("GetPostChats")
	var chts []Chat
	db := GetDB()
	result := db.Where("post_id = ? AND user_id2 = ?", post_id, user_id).Order("id").Limit(to - fr + 1).Offset(fr - 1).Find(&chts)
	if result.Error != nil {
		return nil, result.Error
	}
	return chts, nil
}

func DelChat(cht *Chat) error {
	log.Println("DelChat")
	var dbcht Chat
	db := GetDB()
	result := db.First(&dbcht, cht.ID)
	if result.Error != nil {
		return result.Error
	}
	if cht.SendID != dbcht.SendID {
		return errors.New("permission denied")
	}
	result = db.Delete(&Chat{}, cht.ID)
	return result.Error
}

func GetChatNumber(post_id uint, user_id uint) (int64, error) {
	log.Println("GetChatNumber")
	db := GetDB()
	var count int64
	db.Model(&Chat{}).Where("post_id = ? AND user_id2 = ?", post_id, user_id).Count(&count)
	return count, nil
}

func GetUnseen(user_id uint) (Unseen, error) {
	log.Println("GetUnseen")
	db := GetDB()
	var us Unseen
	result := db.Where("user_id = ?", user_id).First(&us)
	if result.Error != nil {
		return Unseen{}, result.Error
	}
	return us, nil
}

func init() {
	db := GetDB()
	db.AutoMigrate(&Chat{})
	db.AutoMigrate(&Unseen{})
}
