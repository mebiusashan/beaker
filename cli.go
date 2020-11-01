package beaker

import "fmt"

func RunCli() {
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
		fmt.Println("登录成功")
	} else {
		fmt.Println("登录失败")
	}
}

func printHelp() {
	fmt.Printf("请选择操作，输入编号即可。" +
		"\n[1]  文章列表" +
		"\n[2]  删除文章" +
		"\n[3]  添加文章" +
		"\n[4]  下载文章" +
		"\n[5]  查看单页" +
		"\n[6]  删除单页" +
		"\n[7]  添加单页" +
		"\n[8]  下载单页" +
		"\n[9]  添加tweet" +
		"\n[10] 查看tweet" +
		"\n[11] 删除tweet" +
		"\n[12] 查看分类" +
		"\n[13] 添加分类" +
		"\n[14] 修改分类" +
		"\n[15] 删除分类" +
		"\n[16] 查看选项" +
		"\n[17] 清除缓存" +
		"\n请选择操作:")

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
		fmt.Println("选择无效")
	}
}
