package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/zihang5127/easy-operation/command"
	_ "github.com/zihang5127/easy-operation/router"
)


func main() {

	command.InitDatabase();
	command.InitLogger();
	command.InitModel();
	command.InitCommand();

	command.Run()
}