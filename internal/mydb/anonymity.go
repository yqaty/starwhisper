package mydb

import (
	//"cmd/internal/mydb"
	//"cmd/internal/mydb/tools"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type Anonymity struct {
	gorm.Model
	UserId    uint   `json:"user_id"`
	PostId    uint   `json:"post_id"`
	AnonyId   uint   `json:"anony_id"`
	Anonymity string `json:"anonymity"`
}

//if exist, return it, else build it
func GetAnonymity(usrid uint, postid uint) (string, error) {
	log.Println("GetAnonymity")
	db := GetDB()
	var anmy Anonymity
	result := db.Where("user_id = ? AND post_id = ?", usrid, postid).First(&anmy)
	if result.RowsAffected != 0 {
		return anmy.Anonymity, nil
	}
	var aid uint
	for {
		aid = uint(rand.Intn(len(anonymities)) + 1)
		result := db.Where("anony_id = ?", aid).Find(&Anonymity{})
		log.Printf("result.RowsAffected:%d\n", result.RowsAffected)
		if result.RowsAffected == 0 {
			break
		}
	}
	anmy = Anonymity{UserId: usrid, PostId: postid, AnonyId: aid, Anonymity: anonymities[aid]}
	db.Create(&anmy)
	return anmy.Anonymity, nil
}

var anonymities []string

func init() {
	rand.Seed(time.Now().Unix())
	db := GetDB()
	db.AutoMigrate(&Anonymity{})
	anonymities = []string{"水上由岐", "高岛柘榴", "若槻镜", "若槻司", "间宫卓司", "音无彩名", "橘希实香", "间宫羽咲", "悠木皆守", "游行寺夜子", "月社妃", "日向彼方", "伏见理央", "游行寺汀", "曾根美雪", "向日葵", "菜菜木爱丽丝", "灰树由贵", "凛藤华爱美", "蓄井染", "葵矢枢", "森田贤一", "三广幸", "大音灯花", "大音京子", "日向夏咲", "樋口璃璃子", "法月将臣", "期招来那由太", "姬市天美", "芙卡·玛丽尼特", "玖冢蓟子", "玖冢杜鹃子", "寻里耶犹犹", "由芙院御伽", "羽濑久次来", "四十九筮", "高梨子小鸟", "木叶真顷", "安乐村晦", "黑蝶沼志依", "叶深霍", "彬白夜夜萌", "葛笼井桦音", "朝武芳乃", "常陆茉子", "丛雨", "蕾娜·列支敦瑙尔", "鞍马小春", "马庭芦花"}
}
