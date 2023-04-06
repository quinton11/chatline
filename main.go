package main

import (
	"github.com/quinton11/chatline/cmd"
	"github.com/quinton11/chatline/config"
)

func main() {
	//env
	config.LoadConfig()
	cmd.Execute()
}
