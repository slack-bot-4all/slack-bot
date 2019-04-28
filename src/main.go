package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/slack-bot-4all/slack-bot/src/core"
)

func main() {
	core.PrintLogoOnConsole()

	core.Start()
}
