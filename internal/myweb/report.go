package myweb

import (
	"cmd/internal/mydb"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func POSTCommentIDReport(c *gin.Context) {
	session := sessions.Default(c)
	user_id := session.Get("user_id").(uint)
	comment_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	_, err = mydb.GetCommentID(uint(comment_id))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var rpt mydb.Report
	err = c.BindJSON(&rpt)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	rpt.UserID = user_id
	rpt.ContextID = uint(comment_id)
	rpt.IsPost = false
	err = mydb.AddReport(&rpt)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, "Report successfully!")
}

func POSTPostIDReport(c *gin.Context) {
	session := sessions.Default(c)
	user_id := session.Get("user_id")
	post_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	_, err = mydb.GetPost(uint(post_id))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var rpt mydb.Report
	err = c.BindJSON(&rpt)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	rpt.UserID = user_id.(uint)
	rpt.ContextID = uint(post_id)
	rpt.IsPost = true
	err = mydb.AddReport(&rpt)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, "Report successfully!")
}

func InitReportRouter(router *gin.Engine) {
	router.POST("/comment/:id/report", POSTCommentIDReport)
	router.POST("/post/:id/report", POSTPostIDReport)
}
