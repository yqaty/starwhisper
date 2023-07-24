package myweb

import (
	"cmd/internal/mydb"
	"strconv"

	//"cmd/internal/myweb"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type vercode struct {
	Code string `json:"code"`
}

const PageContainUsers = 3

func GETUserRegister(c *gin.Context) {
	log.Println("GETUsersRegister")
	var usr mydb.User
	err := c.BindJSON(&usr)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Println(usr.String())
	rand.Seed(time.Now().Unix())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	err = mydb.SetCode(usr.Email, code)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Println("send email to" + usr.Email)
	err = SendEmail(usr.Email, code)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, "The code has been sent.")
}

func POSTUserRegister(c *gin.Context) {
	log.Println("POSTUsersRegister")
	var usr mydb.User
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	err = json.Unmarshal(bs, &usr)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var code1 vercode
	json.Unmarshal(bs, &code1)
	/*err := c.BindJSON(&usr)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var code1 vercode
	c.BindJSON(&code1)*/
	var code2 string
	code2, err = mydb.GetCode(usr.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("code1:%s,code2:%s\n", code1.Code, code2)
	if code1.Code == code2 {
		err := mydb.AddUser(&usr)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusCreated, usr)
	} else {
		c.String(http.StatusForbidden, "Your code is wrong!")
	}
}

func GETUserLogin(c *gin.Context) {
	log.Println("GETUsersLogin")
	var usr mydb.User
	err := c.BindJSON(&usr)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	op, err := mydb.CheckPassword(&usr)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if op {
		session := sessions.Default(c)
		log.Println("usr.ID", usr.ID)
		session.Set("user_id", usr.ID)
		session.Save()
		c.JSON(http.StatusOK, usr)
	} else {
		c.String(http.StatusForbidden, "Your password is wrong!")
	}
}

func POSTUserForgotPassword(c *gin.Context) {
	log.Println("POSTUsersForgotPassword")
	var usr mydb.User
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	err = json.Unmarshal(bs, &usr)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var code1 vercode
	json.Unmarshal(bs, &code1)
	code2, err := mydb.GetCode(usr.Email)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("code1:%s,code2:%s\n", code1.Code, code2)
	if code1.Code == code2 {
		err = mydb.UpdatePassword(&usr)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.JSON(http.StatusCreated, usr)
	} else {
		c.String(http.StatusBadRequest, "Your code is wrong!")
	}
}

func GETUser(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	users, err := mydb.GetUsers((page-1)*PageContainUsers+1, page*PageContainUsers)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var data []map[string]interface{}
	for i := range users {
		res, err := Type2Map(users[i])
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
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
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if uint(id) != user_id.(uint) {
		c.String(http.StatusBadRequest, "permission denied!")
		return
	}
	posts, err := mydb.UserPosts(uint(id), (page-1)*PageMaxPosts+1, page*PageMaxPosts)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, posts)
}

func GETUserComment(c *gin.Context) {
	session := sessions.Default(c)
	user_id := session.Get("user_id")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if uint(id) != user_id.(uint) {
		c.String(http.StatusBadRequest, "permission denied!")
		return
	}
	cmts, err := mydb.UserComments(uint(id), (page-1)*PageMaxPosts+1, page*PageMaxPosts)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, cmts)
}

func GETUserID(c *gin.Context) {
	session := sessions.Default(c)
	user_id := session.Get("user_id").(uint)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if user_id != uint(id) {
		c.String(http.StatusBadRequest, "permission denied")
		return
	}
	usr, err := mydb.GetUserID(user_id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, usr)
}
func InitUserRouter(router *gin.Engine) {
	log.Println("InitUserRouter")
	sto := cookie.NewStore([]byte("1919810"))
	router.Use(sessions.Sessions("mysession", sto))
	router.GET("/user/register", GETUserRegister)
	router.POST("/user/register", POSTUserRegister)
	router.GET("/user/login", GETUserLogin)
	router.GET("/user/forgot_password", GETUserRegister)
	router.POST("/user/forgot_password", POSTUserForgotPassword)
	router.GET("/article/:id", GETPostID)
	router.GET("/user", AuthRequired(), GETUser)
	router.GET("/user/:id", AuthRequired(), GETUserID)
	router.GET("/user/:id/post", AuthRequired(), GETUserPost)
	router.GET("/user/:id/comment", AuthRequired(), GETUserComment)
}
