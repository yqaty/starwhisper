package myweb

import (
	"cmd/internal/mydb"
	//"cmd/internal/myweb"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GETCommentID(c *gin.Context) {
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
}

func GETComment(c *gin.Context) {
	post_id, err := strconv.Atoi(c.Query("post"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	comment_id, err := strconv.Atoi(c.DefaultQuery("comment", "0"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	floor, err := strconv.Atoi(c.Query("floor"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	cmt, err := mydb.GetComment(uint(post_id), uint(comment_id), uint(floor))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, cmt)
}

func POSTComment(c *gin.Context) {
	session := sessions.Default(c)
	user_id := session.Get("user_id")
	post_id, err := strconv.Atoi(c.Query("post"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	comment_id, err := strconv.Atoi(c.DefaultQuery("comment", "0"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var cmt mydb.Comment
	err = c.BindJSON(&cmt)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	cmt.UserID = user_id.(uint)
	cmt.PostID = uint(post_id)
	cmt.BelongID = uint(comment_id)
	err = mydb.AddComment(&cmt)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, cmt)
}

func DELETECommentID(c *gin.Context) {
	session := sessions.Default(c)
	user_id := session.Get("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var cmt mydb.Comment
	cmt.ID = uint(id)
	cmt.UserID = user_id.(uint)
	err = mydb.DelComment(&cmt)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, "Delete successfully!")
}

func InitCommentRouter(router *gin.Engine) {
	router.GET("/comment/:id", GETCommentID)
	router.GET("/comment", GETComment)
	router.POST("/comment", POSTComment)
	router.DELETE("/comment/:id", DELETECommentID)
}
