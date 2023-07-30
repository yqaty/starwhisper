# 星语 API 文档

## 用户相关

### [POST] `/u/send_code`

给指定邮箱发送验证码

#### 请求

```json
{
	"email":"114514@qq.com"
}
```

#### 响应

```json
{
  "success": true,
  "error": "",
  "data": "The code has been sent successfully!"
}
```

### [POST] `/u/register`

注册用户

#### 请求

```json
{
	"email":"114514@qq.com",
	"password":"1919810",
	"username":"114514",
	"verifycode":"114514"
}
```

#### 响应

```json
{
  "success": true,
  "error": "",
  "data": {
    "id": 4,
    "created_at": "2023-07-29T19:25:50.864003336+08:00",
    "username": "114514",
    "email": "114514@qq.com"
  }
}
```

### [POST] `/u/login`

用户登陆

#### 请求

```json
{
	"email":"114514@qq.com",
	"password:":"1919810"
}
```

#### 响应

```json
{
  "success": true,
  "error": "",
  "data": {
    "id": 1,
    "created_at": "2023-07-29T18:42:52.5121+08:00",
    "username": "114514",
    "email": "114514@qq.com",
    "token": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1..."
  }
}
```

## 星语相关

### [POST] `/p`

发布一条星语

#### 请求

```json
{
  "content":"homo is everywhere!",
  "tags":["倾诉烦恼",
    "悄悄话"]
}
```

#### 响应

```json
{
  "success": true,
  "error": "",
  "data": {
    "id": 1,
    "created_at": "2023-07-29T19:39:53.980897594+08:00",
    "user_id": 1,
    "username": "114514",
    "content": "homo is everywhere!",
    "tags": [
      "倾诉烦恼",
      "悄悄话"
    ]
  }
}
```
### [POST] `/randp`

随机获得一个满足 tag 条件的星语

#### 请求

```json
{
  "tags":["悄悄话"]
}
```

#### 响应

```json
{
  "success": true,
  "error": "",
  "data": {
    "id": 1,
    "created_at": "2023-07-29T19:39:53.980897+08:00",
    "user_id": 1,
    "username": "114514",
    "title": "",
    "content": "homo is everywhere!",
    "tags": [
      "倾诉烦恼",
      "悄悄话"
    ]
  }
}
```

### [GET] `/p/:id`

请求编号为 id 的星语的内容

#### 响应

```json
{
  "success": true,
  "error": "",
  "data": {
    "id": 1,
    "created_at": "2023-07-29T19:39:53.980897+08:00",
    "user_id": 4,
    "username": "www",
    "title": "",
    "content": "homo is everywhere!",
    "tags": [
      "倾诉烦恼",
      "悄悄话"
    ]
  }
}
```

### [DELETE] `/p/:id`

删除编号为 id 的星语

#### 响应

```json
{
  "success": true,
  "error": "",
  "data": "Delete successfully!"
}
```

### [GET] `/pnum`

查看当前用户的星语数量

#### 响应

```json
{
  "success": true,
  "error": "",
  "data": 1 
}
```

### [POST] `/u/p`

倒序获得当前用户的星语

#### 请求

| 参数 | 描述                          |
|------|-------------------------------|
| page | 默认为 1，每页 20 个          |

#### 响应

```json
{
  "success": true,
  "error": "",
  "data": [
    {
      "id": 3,
      "created_at": "2023-07-29T19:50:45.632672+08:00",
      "user_id": 4,
      "username": "www",
      "title": "",
      "content": "1234",
      "tags": []
    },
    {
      "id": 1,
      "created_at": "2023-07-29T19:39:53.980897+08:00",
      "user_id": 4,
      "username": "www",
      "title": "",
      "content": "homo is everywhere!",
      "tags": [
        "倾诉烦恼",
        "悄悄话"
      ]
    }
  ]
}
```

## 通信相关

### [POST] `/chat`

#### 请求

| 参数 | 描述                      |
|------|---------------------------|
| post | 通信的帖子编号            |
| user1| 帖子所有者编号            |
| user2| 另一位通信者的编号        |

```json
{
	"content":"114514"
}
```

#### 响应

```json
{
  "success": true,
  "error": "",
  "data": {
    "id": 1,
    "post_id": 1,
    "user_id1": 1,
    "user_id2": 2,
    "send_id": 2,
    "send_name": "www",
    "created_at": "2023-07-29T22:43:13.508999331+08:00",
    "content": "uh? Are you homo?"
  }
}
```

### [GET] `/chat`

获得通信列表

#### 请求


| 参数 | 描述                      |
|------|---------------------------|
| post | 通信的帖子编号            |
| user1| 帖子所有者编号            |
| user2| 另一位通信者的编号        |
| page | 默认为 1，每页 20 个      |

#### 响应

```json
{
  "success": true,
  "error": "",
  "data": [
    {
      "id": 1,
      "post_id": 1,
      "user_id1": 1,
      "user_id2": 2,
      "send_id": 2,
      "send_name": "www",
      "created_at": "2023-07-29T22:13:32.30867+08:00",
      "content": "uh? Are you homo?"
    },
    {
      "id": 2,
      "post_id": 1,
      "user_id1": 1,
      "user_id2": 2,
      "send_id": 2,
      "send_name": "www",
      "created_at": "2023-07-29T22:13:54.162446+08:00",
      "content": "uh? Are you homo?"
    }
  ]
}
```
### [GET] `/chatnum`

获得通信的总消息数

#### 请求


| 参数 | 描述                      |
|------|---------------------------|
| post | 通信的帖子编号            |
| user1| 帖子所有者编号            |
| user2| 另一位通信者的编号        |

#### 响应

```json
{
  "success": true,
  "error": "",
  "data": 2
}
```

### [DELETE] `/chat/:id`

删除编号为 id 的评论

#### 响应

```json
{
  "success": true,
  "error": "",
  "data": "Delete successfully!"
}
```


### [GET] `/unseen`

返回当前一个用户未读的通信的相关信息

#### 响应

```json
{
  "success": true,
  "error": "",
  "data": {
    "user_id": 4,
    "post_id": 1,
    "send_id": 2
  }
}

```

分别为：当前用户的 id，通信的帖子 id，另一位在该通信中发消息的用户 id。
