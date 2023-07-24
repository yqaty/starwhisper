package mydb

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _db *gorm.DB

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/web")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("%w", err))
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", viper.GetString("postgres.host"), viper.GetString("postgres.user"), viper.GetString("postgres.password"), viper.GetString("postgres.dbname"), viper.GetString("postgres.port"), viper.GetString("postgres.sslmode"), viper.GetString("postgres.TimeZone"))
	//dsn := "host=postgres user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	//dsn := "host=localhost user=yqaty password=yqaty dbname=yqaty port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	time.Sleep(time.Duration(15) * time.Second)
	fmt.Println(dsn)
	_db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	postgresDB, _ := _db.DB()

	postgresDB.SetMaxOpenConns(100)
	postgresDB.SetMaxIdleConns(20)
}

func GetDB() *gorm.DB {
	return _db
}
