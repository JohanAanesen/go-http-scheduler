package main

import (
	"github.com/JohanAanesen/go-http-scheduler/db"
	"github.com/JohanAanesen/go-http-scheduler/model"
)

func main() {
	// Initialize config
	var cfg = Initialize()

	// Initialize timers
	model.InitializeTimers()

	defer db.CloseDB()

	// Run Server
	Run(LoadHTTP(), LoadHTTPS(), cfg.Server)
}
