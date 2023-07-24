package mydb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var redisUrl string

func SetCode(em string, code string) error {
	log.Println("SetCode")
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		return err
	}
	rdb := redis.NewClient(opt)
	res := rdb.Get(context.TODO(), em)
	if res.Val() != "" {
		return errors.New("the verification code has been sent within three minutes")
	}
	log.Println([]byte(code))
	err = rdb.Set(context.TODO(), em, code, 180*time.Second).Err()
	return err
}

func GetCode(em string) (string, error) {
	log.Println("GetCode")
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		return "", err
	}
	rdb := redis.NewClient(opt)
	return rdb.Get(context.TODO(), em).Result()
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("%w", err))
	}
	redisUrl = fmt.Sprintf("redis://%s:%s@%s:%s/%s", viper.GetString("redis.user"), viper.GetString("redis.password"), viper.GetString("redis.host"), viper.GetString("redis.port"), viper.GetString("redis.dbname"))
}
