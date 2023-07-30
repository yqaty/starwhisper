package myweb

import (
	"cmd/internal/mydb"

	//"cmd/internal/myweb"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type vercode struct {
	Code string `json:"verifycode"`
}

type userwithtoken struct {
	mydb.User
	Token string `json:"token"`
}

type returntype struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

const PageContainUsers = 3

func POSTUSendCode(c *gin.Context) {
	log.Println("GETUSendCode")
	var usr mydb.User
	err := c.BindJSON(&usr)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	rand.Seed(time.Now().Unix())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	err = mydb.SetCode(usr.Email, code)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, "failed to set verification code", nil})
		return
	}
	log.Println("send email to" + usr.Email)
	err = SendEmail(usr.Email, code, "【星语】验证码", "您正在使用邮箱注册，"+"验证码："+code+"，切勿将验证码泄漏于他人，本条验证码有效期 5 分钟。")
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, "email sending failed", nil})
		return
	}
	c.JSON(http.StatusOK, returntype{true, "", "The code has been sent successfully!"})
}

func POSTURegister(c *gin.Context) {
	log.Println("POSTURegister")
	var usr mydb.User
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	err = json.Unmarshal(bs, &usr)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	var code1 vercode
	json.Unmarshal(bs, &code1)
	var code2 string
	code2, err = mydb.GetCode(usr.Email)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	log.Printf("code1:%s,code2:%s\n", code1.Code, code2)
	if code1.Code == code2 {
		err := mydb.AddUser(&usr)
		if err != nil {
			c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
			return
		}
		usr.Password = ""
		c.JSON(http.StatusOK, returntype{true, "", usr})
	} else {
		c.JSON(http.StatusOK, returntype{false, "the verifycode is wrong", nil})
	}
}

func POSTUForgotPassword(c *gin.Context) {
	log.Println("POSTUForgotPassword")
	var usr mydb.User
	err := c.BindJSON(&usr)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	rand.Seed(time.Now().Unix())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	err = mydb.SetCode(usr.Email, code)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, "failed to set verification code", nil})
		return
	}
	log.Println("send email to" + usr.Email)
	err = SendEmail(usr.Email, code, "【星语】验证码", "您正在使用邮箱重置密码，"+"验证码："+code+"，切勿将验证码泄漏于他人，本条验证码有效期 5 分钟。")
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, "email sending failed", nil})
		return
	}
	c.JSON(http.StatusOK, returntype{true, "", "The code has been sent successfully!"})
}

func POSTULogin(c *gin.Context) {
	log.Println("POSTULogin")
	var usr mydb.User
	err := c.BindJSON(&usr)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	op, err := mydb.CheckPassword(&usr)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	if op {
		str, err := GenToken(usr.ID, usr.UserName)
		if err != nil {
			c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
			return
		}
		c.Header("Authorization", "Bearer "+str)
		usr.Password = ""
		usrwithtoken := userwithtoken{usr, "Bearer " + str}
		c.JSON(http.StatusOK, returntype{true, "", usrwithtoken})
	} else {
		c.JSON(http.StatusOK, returntype{false, "the password is wrong", nil})
	}
}

func POSTUResetPassword(c *gin.Context) {
	log.Println("POSTUResetPassword")
	var usr mydb.User
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	err = json.Unmarshal(bs, &usr)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	var code1 vercode
	json.Unmarshal(bs, &code1)
	code2, err := mydb.GetCode(usr.Email)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	log.Printf("code1:%s,code2:%s\n", code1.Code, code2)
	if code1.Code == code2 {
		err = mydb.UpdatePassword(&usr)
		if err != nil {
			c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
			return
		}
		usr.Password = ""
		c.JSON(http.StatusOK, returntype{true, "", usr})
	} else {
		c.JSON(http.StatusOK, returntype{false, "the verifycode is wrong", nil})
	}
}

/*func GETUser(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	users, err := mydb.GetUsers((page-1)*PageContainUsers+1, page*PageContainUsers)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	var data []map[string]interface{}
	for i := range users {
		res, err := Type2Map(users[i])
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
		delete(res, "password")
		data = append(data, res)
	}
	c.JSON(http.StatusOK, data)
}

func GETUserPost(c *gin.Context) {
	session := sessions.Default(c)
	user_id := session.Get("user_id")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if uint(id) != user_id.(uint) {
		c.String(http.StatusOK, "permission denied!")
		return
	}
	posts, err := mydb.UserPosts(uint(id), (page-1)*PageMaxPosts+1, page*PageMaxPosts)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.JSON(http.StatusOK, posts)
}

func GETUserComment(c *gin.Context) {
	session := sessions.Default(c)
	user_id := session.Get("user_id")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if uint(id) != user_id.(uint) {
		c.String(http.StatusOK, "permission denied!")
		return
	}
	cmts, err := mydb.UserComments(uint(id), (page-1)*PageMaxPosts+1, page*PageMaxPosts)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.JSON(http.StatusOK, cmts)
}

func GETUserID(c *gin.Context) {
	session := sessions.Default(c)
	user_id := session.Get("user_id").(uint)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if user_id != uint(id) {
		c.String(http.StatusOK, "permission denied")
		return
	}
	usr, err := mydb.GetUserID(user_id)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.JSON(http.StatusOK, usr)
}*/

func InitUserRouter(router *gin.Engine) {
	log.Println("InitUserRouter")
	router.POST("/u/send_code", POSTUSendCode)
	router.POST("/u/register", POSTURegister)
	router.POST("/u/login", POSTULogin)
	router.POST("/u/forgot_password", POSTUForgotPassword)
	router.POST("/u/reset_password", POSTUResetPassword)
	//router.GET("/article/:id", GETPostID)
	//router.GET("/user", AuthRequired(), GETUser)
	//router.GET("/user/:id", AuthRequired(), GETUserID)
	//router.GET("/user/:id/post", AuthRequired(), GETUserPost)
	//router.GET("/user/:id/comment", AuthRequired(), GETUserComment)
}
