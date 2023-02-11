package main

import (
	"cc/app/internal/app"
	"math/rand"
	"time"
)

func init() {
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc
	
	rand.Seed(time.Now().UnixNano())
}

func main() {
	app.New().
		Run()
}
