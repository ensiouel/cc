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
