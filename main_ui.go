//go:build ui

package main

import (
	"embed"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/issueye/version-mana/docs"
	"github.com/issueye/version-mana/internal/global"
	"github.com/issueye/version-mana/internal/initialize"
)

// @title       版本管理服务
// @version     V0.1
// @description 版本管理服务

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization

//go:embed assets/*
var ui embed.FS

func main() {
	global.TagName = "UI"
	global.Assets = ui

	initialize.Initialize()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
