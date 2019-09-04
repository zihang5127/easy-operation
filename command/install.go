package command

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/zihang5127/easy-operation/model"
	"os"
	"strings"
)

// Install 安装
//使用方式：easy-operation install -username=admin -password=123456 -email=512796048@qq.com
func Install() {

	if len(os.Args) > 2 && os.Args[1] == "install" {

		username := flag.String("username", "username", "Administrator username.")
		pwd := flag.String("password", "", "Administrator password.")
		email := flag.String("email", "", "Administrator email.")

		_ = flag.CommandLine.Parse(os.Args[2:])

		password := strings.TrimSpace(*pwd)

		if password == "" {
			fmt.Println("Administrator password  is required.")
			os.Exit(0)
		}
		if *email == "" {
			fmt.Println("Administrator email is required")
			os.Exit(0)
		}

		user := model.NewUser()
		user.Username = *username
		user.Password = password
		user.Email = *email

		if err := user.Add(); err != nil {
			logs.Error("%s: ", err)
			os.Exit(0)
		}
		os.Exit(0)
	}
}
