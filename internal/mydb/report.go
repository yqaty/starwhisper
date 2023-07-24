package mydb

import (
	//"cmd/internal/mydb"
	//"cmd/internal/mydb/tools"
	"errors"
	"log"

	"gorm.io/gorm"
)

type Report struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	ContextID uint   `json:"context_id"`
	IsPost    bool   `json:"is_post"`
	Reason    string `json:"reason"`
	State     string `gorm:"default:dealing" json:"state"`
}

func AddReport(rpt *Report) error {
	log.Println("AddReport")
	db := GetDB()
	result := db.Create(rpt)
	return result.Error
}

func GetReport(rptid uint) (Report, error) {
	log.Println("Getreport")
	var rpt Report
	rpt.ID = rptid
	db := GetDB()
	result := db.First(&rpt)
	if result.RowsAffected == 0 {
		return Report{}, errors.New("the report does not exist")
	}
	return rpt, nil
}

//descending order
func UserReport(usrid uint) ([]Report, error) {
	log.Println("UserReport")
	var rpts []Report
	db := GetDB()
	result := db.Where("user_id = ?", usrid).Order("id desc").Find(&rpts)
	if result.Error != nil {
		return nil, result.Error
	}
	return rpts, nil
}

func init() {
	db := GetDB()
	db.AutoMigrate(&Report{})
}
