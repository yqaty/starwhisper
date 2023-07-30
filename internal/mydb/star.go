package mydb

import "log"

//"cmd/internal/mydb/tools"

type Star struct {
	PostID uint `gorm:"primaryKey"`
	UserID uint `gorm:"primaryKey"`
}

func IsStar(star *Star) (bool, error) {
	log.Println("IsStar")
	db := GetDB()
	result := db.Where("post_id = ? AND user_id = ?", star.PostID, star.UserID).First(star)
	if result.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func AddStar(star *Star) error {
	log.Println("AddStar")
	db := GetDB()
	result := db.Save(star)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DelStar(star *Star) error {
	log.Println("DelStar")
	db := GetDB()
	db.Delete(star)
	return nil
}

func GetStar(user_id uint, fr int, to int) ([]uint, error) {
	log.Println("GetStar")
	db := GetDB()
	var stars []Star
	result := db.Where("user_id = ?", user_id).Limit(to - fr + 1).Offset(fr - 1).Find(&stars)
	if result.Error != nil {
		return nil, result.Error
	}
	post_ids := make([]uint, len(stars))
	for i, v := range stars {
		post_ids[i] = v.PostID
	}
	return post_ids, nil
}

func GetUserStarNumber(user_id uint) (int64, error) {
	log.Println("GetUserStarNumber")
	db := GetDB()
	var count int64
	db.Model(&Star{}).Where("user_id = ?", user_id).Count(&count)
	return count, nil
}

func init() {
	db := GetDB()
	db.AutoMigrate(&Star{})
}
