package mydb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var _opt *redis.Options

func SetCode(em string, code string) error {
	log.Println("SetCode")
	rdb := redis.NewClient(_opt)
	/*res := rdb.Get(context.TODO(), em)
	if res.Val() != "" {
		return errors.New("the verification code has been sent within five minutes")
	}
	log.Println([]byte(code))*/
	err := rdb.Set(context.TODO(), em, code, 300*time.Second).Err()
	return err
}

func GetCode(em string) (string, error) {
	log.Println("GetCode")
	rdb := redis.NewClient(_opt)
	return rdb.Get(context.TODO(), em).Result()
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("%w", err))
	}
	redisUrl := fmt.Sprintf("redis://%s:%s@%s:%s/%s", viper.GetString("redis.user"), viper.GetString("redis.password"), viper.GetString("redis.host"), viper.GetString("redis.port"), viper.GetString("redis.dbname"))
	_opt, err = redis.ParseURL(redisUrl)
	if err != nil {
		panic(fmt.Errorf("%w", err))
	}
}
