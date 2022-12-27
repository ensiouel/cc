package main

import (
	"cc/app/internal/app"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	app.New().
		Run()
}

/*

* Shorten
GET /api/shortens/:id
{
	"id": "google",
	"title": "Google",
	"long_url": "https://google.com",
	"short_url": "blbm.cc/google",
	"created_at": "2022-12-22T22:39:25+03:00"
}

POST /api/shortens
REQUEST:
{
	"id": "google",
	"long_url": "https://google.com",
	"title": "Google"
}
RESPONSE:
{
	"id": "google",
	"title": "Google",
	"long_url": "https://google.com",
	"short_url": "blbm.cc/google",
	"created_at": "2022-12-22T22:39:25+03:00"
}

PATCH /api/shortens/:id
{
	"title": "Yandex"
}

DELETE /api/shortens/:id
GET /api/shortens/:id/[clicks, platforms, referrers, os]
{
	"clicks": [
		{
			"count": 0,
			"date": "2022-12-22T22:39:25+03:00"
		}
	]
}

{
	"platforms": [
		{
			"name": "Windows",
			"date": "2022-12-22T22:39:25+03:00"
			"count": 0
		}
	]
}

* User
GET  /api/users/:id
GET  /api/users/:id/shortens
{
	"response": [
		{
			"id": "dfk5",
			"title": "Google",
			"long_url": "https://google.com",
			"short_url": "blbm.cc/dfk5",
			"created_at": "2022-12-22T22:39:25+03:00"
		},
		{
			"id": "XIuWx",
			"title": "twitch.tv",
			"long_url": "https://twitch.tv",
			"short_url": "blbm.cc/XIuWx",
			"created_at": "2022-12-22T22:40:03+03:00"
		}
	]
}

POST /api/users/signin
POST /api/users/signup

{
	"error": {
		"code": 0,
		"message": ""
	}
}


{
	"clicks": [
		{"count": 0, "date": ""},
		{"count": 0, "date": ""},
		{"count": 0, "date": ""},
	],
	"unit": "day",
	"units": 3
}

{
	"clicks": 100
}

*/
