package myweb

import (
	"cmd/internal/mydb"
	"encoding/json"
	"log"

	//"cmd/internal/myweb"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const PageMaxPosts = 3

const layout = "2006-01-02T15:04:05Z07:00"

func Type2Map(data interface{}) (map[string]interface{}, error) {
	jsn, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var ndata map[string]interface{}
	err = json.Unmarshal(jsn, &ndata)
	if err != nil {
		return nil, err
	}
	return ndata, nil
}

func Posts2Map(posts []interface{}) []map[string]interface{} {
	var mp []map[string]interface{}
	for i := range posts {
		pt, _ := Type2Map(posts[i])
		delete(pt, "user_id")
		mp = append(mp, pt)
	}
	return mp
}

func GETPost(c *gin.Context) {
	fr, err := time.Parse(layout, c.DefaultQuery("from", "0002-01-01T00:00:00Z"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	to, err := time.Parse(layout, c.DefaultQuery("to", "2050-01-01T00:00:00Z"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	fr = fr.Add(-8 * time.Hour)
	to = to.Add(-8 * time.Hour)
	//fr = fr.In(time.FixedZone("UTC", 8*60*60))
	//to = to.In(time.FixedZone("UTC", 8*60*60))
	title := c.DefaultQuery("title", "")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	posts, err := mydb.SearchPosts(title, fr, to, (page-1)*PageMaxPosts+1, page*PageMaxPosts)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	/*var intervals = make([]interface{}, len(posts))
	for i := range posts {
		intervals = append(intervals, posts[i])
	}*/
	c.JSON(http.StatusOK, posts)
}

func POSTPost(c *gin.Context) {
	session := sessions.Default(c)
	user_id := session.Get("user_id")
	log.Println("debug", user_id.(uint))
	var pt mydb.Post
	err := c.BindJSON(&pt)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	pt.UserID = user_id.(uint)
	err = mydb.AddPost(&pt)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, pt)
}

func PUTPostID(c *gin.Context) {
	session := sessions.Default(c)
	user_id := session.Get("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var pt mydb.Post
	err = c.BindJSON(&pt)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	pt.UserID = user_id.(uint)
	pt.ID = uint(id)
	err = mydb.UpdatePost(&pt)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, pt)
}

func DELETEPostID(c *gin.Context) {
	session := sessions.Default(c)
	user_id := session.Get("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var pt mydb.Post
	pt.ID = uint(id)
	pt.UserID = user_id.(uint)
	err = mydb.DelPost(&pt)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, "Delete successfully!")
}

func GETPostID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	pt, err := mydb.GetPost(uint(id))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, pt)
}

func InitPostRouter(router *gin.Engine) {
	router.Use(AuthRequired())
	router.GET("/post", GETPost)
	router.POST("/post", POSTPost)
	router.PUT("/post/:id", PUTPostID)
	router.DELETE("/post/:id", DELETEPostID)
	router.GET("/post/:id", GETPostID)
}
