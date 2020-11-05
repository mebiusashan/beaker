package beaker

import (
	"bufio"
	"fmt"
	"os"
)

func RunCli() {
	server, serr := CMDGetServer()
	if serr != nil {
		fmt.Println("Please input your beaker blog admin url(like https://example.com:9092):")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		server = input.Text()
		e := CMDSetServer(server)
		if e != nil {
			fmt.Println(e)
			return
		}
	}

	HOST = server

	key, err := CMDGetLocalKey()
	if err != nil {
		///本地key没有或者错误，直接到登录
		login()
		return
	}
	checkKey := CMDCheck(key)
	if !checkKey {
		//key错误，重新登录
		login()
		return
	}
	printHelp()
}

func login() {
	k := CMDPing()
	if k == nil {
		return
	}
	sk := CMDLogin(k)
	if sk == nil {
		return
	}
	rel := CMDCheck(sk)
	if rel {
		fmt.Println( GetLanguage("loginSucc"))
	} else {
		fmt.Println( GetLanguage("loginfail"))
	}
}

func printHelp() {
	fmt.Printf(GetLanguage("PleaseSelectOperation"))

	var action string
	_, _ = fmt.Scanln(&action)
	switch action {
	case "1":
		CMDArcAll()
	case "2":
		CMDArcDel()
	case "3":
		CMDArcAdd()
	case "4":
		CMDArcDown()
	case "5":
		CMDPagAll()
	case "6":
		CMDPagDel()
	case "7":
		CMDPagAdd()
	case "8":
		CMDPagDown()
	case "9":
		CMDTweAdd()
	case "10":
		CMDTweAll()
	case "11":
		CMDTweDel()
	case "12":
		CMDCatAll()
	case "13":
		CMDCatAdd()
	case "14":
		CMDCatEdit()
	case "15":
		CMDCatDel()
	case "16":
		CMDSeeOpts()
	case "17":
		CMDClearAllCache()
	default:
		fmt.Println(GetLanguage("InvalidInstruction"))
	}
}
