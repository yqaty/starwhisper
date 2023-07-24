package myweb

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const apistr = `
仅支持JSON格式www
GET /api/v1 获取api
GET /user/register 获取邮箱验证码，json需有 email
POST /user/register 验证并注册，json需有 email，password，code，可选user_name，gender
GET /user/login 登录，json需有 email，password，响应报文会给 session
GET /user/forgot_password 忘记密码时获取邮箱验证码，json需有 email
POST /user/forgot_password 验证并重置，json需有 email，password，code
GET /article/:id 帖子的外部链接

以下均需身份验证

GET /user?page=x 查看用户信息，每页3个，page 默认为1 
GET /user/:id 查看某个用户的信息
GET /user/:id/post?page=1 查看某个用户的帖子，每页3个，page 默认为1
GET /user/:id/comment?page=1 查看某个用户的评论，每页3个
GET /post?from=2023-07-22T05:20:00Z&to=2023-08-22T05:20:00Z&title=xxx&page=1 搜索标题包含 title，创建时间从 from 到 to 的帖子，每页3个，page 默认为1，fr，to可省略
POST /post 发表帖子，json需 title，context
PUT /post/:id 更新帖子，json需有 title，context
DELETE /post/:id 删除某个帖子
GET /post/:id 查看某个帖子
GET /comment/:id 查看某个评论
GET /comment?post=1&comment=1&floor=1 查看主帖为 post 的评论，comment默认为0，为0代表是评论主帖的，不为0表示是评论 comment 的，floor 为楼层数
POST /comment?post=1&comment=1 发表主帖为 post 的评论，comment默认为0，为0表示是评论主帖的，不为0表示是评论 comment 的
DELETE /comment/:id 删除某个评论
POST /post/:id/report 举报某个帖子，json 需有 reason
POST /comment/:id/report 举报某个评论，json 需有 reason

`

func GETApi(c *gin.Context) {
	c.String(http.StatusOK, apistr)
}

func InitApiRouter(router *gin.Engine) {
	router.GET("/api/v1", GETApi)
}
