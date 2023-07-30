package myweb

import (
	"cmd/internal/mydb"
	"encoding/json"
	"log"

	//"cmd/internal/myweb"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const _postn = 20

//const layout = "2006-01-02T15:04:05Z07:00"

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

/*func GETPost(c *gin.Context) {
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
	//var intervals = make([]interface{}, len(posts))
	//for i := range posts {
//		intervals = append(intervals, posts[i])
//	}
	c.JSON(http.StatusOK, posts)
}*/

func GETPost(c *gin.Context) {
	log.Println("GETPost")
	tags := struct {
		Tags []string `json:"tags"`
	}{}
	err := c.BindJSON(&tags)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	log.Println(tags.Tags)
	if len(tags.Tags) == 0 {
		post, err := mydb.GetRandPost()
		if err != nil {
			c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
			return
		}
		tags.Tags, err = mydb.GetPostTags(post.ID)
		if err != nil {
			c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
			return
		}
		c.JSON(http.StatusOK, returntype{true, "", mydb.PostWithTags{Post: post, Tags: tags.Tags}})
	} else {
		post_id, err := mydb.GetPostByTags(tags.Tags)
		if err != nil {
			c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
			return
		}
		postwithtags, err := mydb.GetPostWithTagByID(post_id)
		if err != nil {
			c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
			return
		}
		c.JSON(http.StatusOK, returntype{true, "", postwithtags})
	}
}

func POSTPost(c *gin.Context) {
	log.Println("POSTPost")
	user_id, _ := c.Get("user_id")
	username, _ := c.Get("username")
	log.Println("debug", user_id.(uint))
	var pt mydb.PostWithTags
	err := c.BindJSON(&pt)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	pt.UserID = user_id.(uint)
	pt.Username = username.(string)
	err = mydb.AddPost(&pt.Post)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	err = mydb.AddTags(pt.ID, pt.Tags)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, returntype{true, "", pt})
}

func DELETEPostID(c *gin.Context) {
	log.Println("DELETEPostID")
	user_id, _ := c.Get("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	var pt mydb.Post
	pt.ID = uint(id)
	pt.UserID = user_id.(uint)
	err = mydb.DelPost(&pt)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	mydb.DelPostTags(uint(id))
	c.JSON(http.StatusOK, returntype{true, "", "Delete successfully!"})
}

func GETPostNumber(c *gin.Context) {
	log.Println("GETPostNumbers")
	user_id, _ := c.Get("user_id")
	count, err := mydb.GetUserPostNumber(user_id.(uint))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, returntype{true, "", count})
}

func GETPostID(c *gin.Context) {
	log.Println("GetPostID")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	pt, err := mydb.GetPostWithTagByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, returntype{true, "", pt})
}

func POSTStarPost(c *gin.Context) {
	log.Println("POSTPostStar")
	user_id, _ := c.Get("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	err = mydb.AddStar(&mydb.Star{PostID: uint(id), UserID: user_id.(uint)})
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, returntype{true, "", "Star successfully!"})
}

func DELETEStarPost(c *gin.Context) {
	log.Println("DELETEStarPost")
	user_id, _ := c.Get("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	err = mydb.DelStar(&mydb.Star{PostID: uint(id), UserID: user_id.(uint)})
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, returntype{true, "", "DELETE successfully!"})
}

func GETStarPnum(c *gin.Context) {
	log.Println("GETStarPnum")
	user_id, _ := c.Get("user_id")
	count, err := mydb.GetUserStarNumber(user_id.(uint))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, returntype{true, "", count})
}

func GETStarPost(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	post_ids, err := mydb.GetStar(user_id.(uint), _postn*(page-1)+1, _postn*page)
	postswithtags := make([]mydb.PostWithTags, len(post_ids))
	for i, id := range post_ids {
		postswithtags[i], _ = mydb.GetPostWithTagByID(id)
	}
	c.JSON(http.StatusOK, returntype{true, err.Error(), postswithtags})
}

func GETUserPost(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	posts, err := mydb.GetUserPosts(user_id.(uint), _postn*(page-1)+1, _postn*page)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	postswithtags, err := mydb.Posts2PostsWithTags(posts)
	if err != nil {
		c.JSON(http.StatusOK, returntype{false, err.Error(), nil})
		return
	}
	c.JSON(http.StatusOK, returntype{true, "", postswithtags})
}

func InitPostRouter(router *gin.Engine) {
	router.POST("/randp", GETPost)
	router.POST("/p", JWTAuthMiddleware(), POSTPost)
	router.GET("/p/:id", JWTAuthMiddleware(), GETPostID)
	router.DELETE("/p/:id", JWTAuthMiddleware(), DELETEPostID)
	router.GET("/pnum", JWTAuthMiddleware(), GETPostNumber)
	router.POST("/star/p/:id", JWTAuthMiddleware(), POSTStarPost)
	router.DELETE("/star/p/:id", JWTAuthMiddleware(), DELETEStarPost)
	router.GET("/star/pnum", JWTAuthMiddleware(), GETStarPnum)
	router.GET("/star/p", JWTAuthMiddleware(), GETStarPost)
	router.GET("/u/p", JWTAuthMiddleware(), GETUserPost)
}
