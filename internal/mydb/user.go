package mydb

import (
	//"cmd/internal/mydb/tools"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
	Password string `json:"password"`
}

func (usr *User) String() string {
	return fmt.Sprintf("ID:%d,UserName:%s,Email:%s,Gender:%s,Password:%s\n", usr.ID, usr.UserName, usr.Email, usr.Gender, usr.Password)
}

func GetUsers(fr int, to int) ([]User, error) {
	log.Println(fr, to)
	if fr > to {
		return nil, errors.New("request illegal")
	}
	var users []User
	db := GetDB()
	result := db.Limit(to - fr + 1).Offset(fr - 1).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func EmailIsExist(email string) bool {
	db := GetDB()
	result := db.Where("email = ?", email).First(&User{})
	return result.RowsAffected > 0
}

func AddUser(newusr *User) error {
	log.Println("AddUser")
	if EmailIsExist(newusr.Email) {
		return errors.New("the email has been registered")
	}
	db := GetDB()
	psd, err := bcrypt.GenerateFromPassword([]byte(newusr.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newusr.Password = string(psd)
	if result := db.Create(&newusr); result.Error != nil {
		return err
	}
	return nil
}

func CheckPassword(usr *User) (bool, error) {
	log.Println("CheckPassword")
	db := GetDB()
	var dbusr User
	if result := db.Where("email = ?", usr.Email).First(&dbusr); result.Error != nil {
		return false, result.Error
	}
	err := bcrypt.CompareHashAndPassword([]byte(dbusr.Password), []byte(usr.Password))
	*usr = dbusr
	return err == nil, err
}

func UpdatePassword(usr *User) error {
	log.Println("UpdatePassword")
	db := GetDB()
	var dbusr User
	if result := db.Where("email = ?", usr.Email).First(&dbusr); result.Error != nil {
		return result.Error
	}
	var psd []byte
	psd, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	dbusr.Password = string(psd)
	db.Save(&dbusr)
	*usr = dbusr
	return nil
}

func GetUserID(id uint) (User, error) {
	db := GetDB()
	var usr User
	result := db.First(&usr, id)
	if result.Error != nil {
		return User{}, result.Error
	}
	return usr, nil
}

func init() {
	db := GetDB()
	db.AutoMigrate(&User{})
}
