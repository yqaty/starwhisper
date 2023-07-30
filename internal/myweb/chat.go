package myweb

import (
	"cmd/internal/mydb"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const _commentn = 20

/*func GETCommentID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	cmt, err := mydb.GetCommentID(uint(id))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, cmt)
}*/

func GETChats(c *gin.Context) {
	log.Println("GETChats")
	user_id, _ := c.Get("user_id")
	post_id, err := strconv.Atoi(c.Query("post"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	user_id1, err := strconv.Atoi(c.Query("user1"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	user_id2, err := strconv.Atoi(c.Query("user2"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	chts, err := mydb.GetPostChats(uint(post_id), uint(user_id2), _commentn*(page-1)+1, _commentn*page)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	mydb.DelUnseen(&mydb.Unseen{UserID: user_id.(uint), PostID: uint(post_id), SendID: uint(user_id1) ^ uint(user_id2) ^ user_id.(uint)})
	c.JSON(http.StatusOK, returntype{true, "", chts})
}

func POSTChat(c *gin.Context) {
	log.Println("POSTChat")
	send_id, _ := c.Get("user_id")
	send_name, _ := c.Get("username")
	post_id, err := strconv.Atoi(c.Query("post"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	user_id1, err := strconv.Atoi(c.Query("user1"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	user_id2, err := strconv.Atoi(c.Query("user2"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	if user_id1 == user_id2 {
		c.JSON(http.StatusOK, returntype{false, "You cannot reply to yourself", nil})
	}
	var cht mydb.Chat
	err = c.BindJSON(&cht)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	cht.UserID1 = uint(user_id1)
	cht.UserID2 = uint(user_id2)
	cht.PostID = uint(post_id)
	cht.SendID = send_id.(uint)
	cht.SendName = send_name.(string)
	err = mydb.AddChat(&cht)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	mydb.AddUnseen(&mydb.Unseen{UserID: cht.UserID1 ^ cht.UserID2 ^ cht.SendID, PostID: cht.PostID, SendID: cht.SendID})
	c.JSON(http.StatusOK, returntype{true, "", cht})
}

func DELETEChatByID(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	cmt_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	var cht mydb.Chat
	cht.SendID = user_id.(uint)
	cht.ID = uint(cmt_id)
	err = mydb.DelChat(&cht)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, returntype{true, "", "Delete successfully!"})
}

func GETChatNumber(c *gin.Context) {
	log.Println("GETChatNumber")
	post_id, err := strconv.Atoi(c.Query("post"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	user_id2, err := strconv.Atoi(c.Query("user2"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	count, _ := mydb.GetChatNumber(uint(post_id), uint(user_id2))
	c.JSON(http.StatusOK, returntype{true, "", count})
}

func GETUnseen(c *gin.Context) {
	log.Println("GETUnseen")
	user_id, _ := c.Get("user_id")
	us, err := mydb.GetUnseen(user_id.(uint))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, returntype{true, "", us})
}

func InitCommentRouter(router *gin.Engine) {
	//router.GET("/comment/:id", GETCommentID)
	router.GET("/unseen", JWTAuthMiddleware(), GETUnseen)
	router.GET("/chatnum", JWTAuthMiddleware(), GETChatNumber)
	router.POST("/chat", JWTAuthMiddleware(), POSTChat)
	router.GET("/chat", JWTAuthMiddleware(), GETChats)
	router.DELETE("/chat/:id", JWTAuthMiddleware(), DELETEChatByID)
	//router.GET("/cnum/:id", JWTAuthMiddleware(), GETCommentNumber)
}
