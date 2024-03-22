# starwhisper

## 一个用 gin，gorm 框架编写的漂流瓶网站

### 使用方法（谔谔）

本项目用 postgresql 储存用户信息帖子评论举报信息，用 redis 储存验证码。

```
{
    "port":"8000",
    "domain":"",

    "postgres":{
        "host":"postgres",
        "user":"postgres",
        "password":"postgres",
        "dbname":"postgres",
        "port":"5432",
        "sslmode":"disable",
        "TimeZone":"Asia/Shanghai"
    },

    "redis":{
        "user":"",
        "password":"",
        "host":"redis",
        "port":"6379",
        "dbname":"0"
    },

	"email":{
		"email":"",
		"server":"",
		"port":"",
		"auth_code":""
	}
}
```
请先配置 config.json , domain 为域名，email 为发邮件的地址，服务商，端口，授权码等信息，使用 SMTP 协议。
然后在命令行输入 ```make```，程序会在 docker 容器中运行。
使用 ``` GET /api/v1 ``` 方法获取 api。

