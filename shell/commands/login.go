package commands

import (
	"fmt"
	pl "github.com/Mrs4s/power-liner"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/shell"
	"github.com/Mrs4s/six-cli/six_cloud"
	"strings"
)

func init() {
	alias["Login"] = []string{}
	explains["Login"] = "登录6Pan账号"
}

func (CommandHandler) Login(c *pl.Context) {
	var (
		username string
		password string
		args     = models.FilterStrings(c.RawArgs, func(s string) bool { return strings.TrimSpace(s) != "" })
	)
	if len(args) == 2 {
		username = args[0]
		password = args[1]
	} else {
		username, _ = shell.App.ReadLine("请输入用户名或手机号: ")
		password, _ = shell.App.ReadPassword("请输入密码: ")
	}
	user, err := six_cloud.LoginWithUsernameOrPhone(username, password)
	if err != nil {
		fmt.Println("[!] 登录失败: " + err.Error())
		return
	}
	shell.CurrentUser = user
	shell.CurrentPath = "/"
	var flag bool
	for _, su := range shell.SavedUsers {
		if su.Identity == user.Identity {
			flag = true
		}
	}
	if !flag {
		models.DefaultConf.QingzhenTokens = append(models.DefaultConf.QingzhenTokens, user.Client.QingzhenToken)
		shell.SavedUsers = append(shell.SavedUsers, user)
	}
	fmt.Println("[+] 登录完成, 欢迎: " + user.Username)
	fmt.Println()
	if len(shell.SavedUsers) > 1 {
		printUserList()
	}
	models.DefaultConf.SaveFile("config.json")
	refreshPrompt()
}
