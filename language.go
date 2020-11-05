package beaker

type Language int

const (
	Language_EN Language = 0
	Language_CN Language = 1
)

var curLanguage Language = Language_EN
var lang map[string]string

func GetLanguage(key string) string{
	if 	lang == nil{
		SetLanguage(Language_EN)
	}
	return lang[key];
}

func SetLanguage(language Language) {
	lang = make(map[string]string)
	if language == Language_CN {
		initLangCN()
		return
	}
	initLangEN();
}

func initLangCN(){
	lang["ArticlAddedSucc"] = "文章添加成功"
	lang["ArticleDelSucc"] = "文章删除成功"
	lang["AllArticle"] = "全部文章"
	lang["Article"] = "文章"
	lang["EnterArticleIDToDel"] = "请输入要删除的文章ID :"
	lang["IDError"] = "ID错误"
	lang["Title"] = "标题"
	lang["mdFileNotFound"] = "没有找到Markdown(.md)文件"
	lang["ClearCacheSucc"] = "清除缓存成功"
	lang["Done"] = "完成"
	lang["name"] = "名称"
	lang["value"] = "值"
	lang["tweetAddSucc"] = "tweet添加成功"
	lang["tweetDelSucc"] = "tweet删除成功"
	lang["Content"] = "内容"
	lang["CreateTime"] = "创建时间"
	lang["DecodingFailed"] = "解码失败"
	lang["username"] = "用户名:"
	lang["password"] = "密码:"
	lang["loginSucc"] = "登录成功"
	lang["loginfail"] = "登录失败"
	lang["InvalidInstruction"] = "无效指令"
	lang["PleaseSelectOperation"] = "请选择操作，输入编号即可。" +
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
	"\n请选择操作:"
	lang["EnterTweetIDToDeleted"] = "请输入要删除的Tweet ID:"
	lang["ConfirmInput"] = "确认输入 (y or n):"
	lang["CategoryID"] = "分类ID"
	lang["EnterTweetContent"] = "请输入Tweet内容:"
	lang["AreYouSureAddTweet"] = "确认添加Tweet么?"
	lang["AreYouSureDelTweet"] = "确认删除Tweet么?"
	lang["PageAddSucc"] = "单页添加成功"
	lang["PageDelSucc"] = "单页删除成功"
	lang["EnterIDPageToDel"] = "请输入要删除的单页ID:"
	lang["AreYouSureDeletePage"] = "确认删除单页么?"
	lang["EnterIDPageToDownloaded"] = "请输入要下载的单页ID:"
	lang["AreYouSureDownloadPage"] = "确认下载单页么?"
	lang["ArticleDownloadedSucceAndIn"] = "文章下载成功，存储到:"
	lang["EnterPageTitle"] = "请输入单页标题:"
	lang["AreYouSureAddetePage"] = "确认添加单页么?"
	lang["AreYouSureAddeteArc"] = "确认添加文章么?"
	lang["EnterArticleTitle:"] = "请输入文章标题:"
	lang["EnterArticleCategoryID"] = "请输入文章分类ID:"
	lang["AreYouSureDeleteArticle"] = "确认删除文章么?"
	lang["EnterIDArticleDownload"] = "请输入要下载的文章ID:"
	lang["AreYouSureDownloadArticle"] = "确认下载文章么?"
	lang["NotNull"] = "不允许出现空值"
	lang["catAddSucc"] = "分类添加成功"
	lang["catDelSucc"] = "分类删除成功"
	lang["catall"] = "所有分类"
	lang["CategoryModifiedSucc"] = "分类修改成功"
	lang["showname"] = "显示名"
	lang["path"] = "Path"
	lang["EnterDisplayName"] = "请输入显示名:"
	lang["EnterPath"] = "请输入路径名:"
	lang["ConfirmModifyCategory"] = "确认修改分类么?"
	lang["AreYouSureAddCategory"] = "确认添加新分类么?"
	lang["EnterCatIDDel"] = "请输入要删除的分类ID:"
	lang["moveCatId"] = "请输入当前分类下文章移动到的分类ID:"
	lang["AreYouSureDelCategory"] = "确认删除分类么?"
	lang["EnterCategoryIDModified"] = "请输入要修改的分类ID:"
}

func initLangEN(){
	lang["ArticlAddedSucc"] = "Article added successfully"
	lang["ArticleDelSucc"] = "Article deleted successfully:";
	lang["AllArticle"] = "All Article";
	lang["Article"] = "Article";
	lang["EnterArticleIDToDel"] = "Please enter the article ID to be deleted :";
	lang["IDError"] = "ID error";
	lang["Title"] = "Title";
	lang["mdFileNotFound"] = "Markdown file not found"
	lang["ClearCacheSucc"] = "Clear cache successfully"
	lang["Done"] = "Done"
	lang["name"] = "Name"
	lang["value"] = "Value"
	lang["tweetAddSucc"] = "Tweet added successfully"
	lang["tweetDelSucc"] = "Tweet deleted successfully"
	lang["Content"] = "Content"
	lang["CreateTime"] = "Add Time"
	lang["DecodingFailed"] = "Decoding failed"
	lang["username"] = "User name"
	lang["password"] = "Password"
	lang["loginSucc"] = "Login successfully"
	lang["loginfail"] = "Login failed"
	lang["InvalidInstruction"] = "Invalid instruction"
	lang["PleaseSelectOperation"] = "Please select an operation and enter the number."+
	"\n[1]  Article list" +
	"\n[2]  Delete article" +
	"\n[3]  Add article" +
	"\n[4]  Download article" +
	"\n[5]  Page list" +
	"\n[6]  Delete Page" +
	"\n[7]  Add Page" +
	"\n[8]  Download Page" +
	"\n[9]  Add Tweet" +
	"\n[10] Tweet list" +
	"\n[11] Delete Tweet" +
	"\n[12] Category list" +
	"\n[13] Add Category" +
	"\n[14] Modify Category" +
	"\n[15] Delete Category" +
	"\n[16] View options" +
	"\n[17] Clear cache" +
	"\nPlease choose an operation:"
	lang["EnterTweetIDToDeleted"] = "Please enter the Tweet ID to be deleted"
	lang["ConfirmInput"] = "Confirm input (y or n):"
	lang["CategoryID"] = "Category ID"
	lang["EnterTweetContent"] = "Please enter Tweet content:"
	lang["AreYouSureAddTweet"] = "Are you sure to add Tweet?"
	lang["AreYouSureDelTweet"] = "Are you sure to delete Tweet?"
	lang["PageAddSucc"] = "Page added successfully"
	lang["PageDelSucc"] = "Page deleted successfully"
	lang["EnterIDPageToDel"] = "Please enter the ID of the single page to be deleted:"
	lang["AreYouSureDeletePage"] = "Are you sure to delete the single page?"
	lang["EnterIDPageToDownloaded"] = "Please enter the ID of the single page to be downloaded:"
	lang["AreYouSureDownloadPage"] = "Are you sure to download the single page?"
	lang["ArticleDownloadedSucceAndIn"] = "The article is downloaded successfully and stored in:"
	lang["EnterPageTitle"] = "Please enter a single page title:"
	lang["AreYouSureAddetePage"] = "Are you sure to add the single page?"
	lang["AreYouSureAddeteArc"] = "Are you sure to add the article?"
	lang["EnterArticleTitle:"] = "Please enter the article title:"
	lang["EnterArticleCategoryID"] = "Please enter the article category ID:"
	lang["AreYouSureDeleteArticle"] = "Are you sure to delete the article?"
	lang["EnterIDArticleDownload"] = "Please enter the ID of the article you want to download:"
	lang["AreYouSureDownloadArticle"] = "Are you sure to download the article?"
	lang["NotNull"] = "Null values ​​are not allowed"
	lang["catAddSucc"] = "Category added successfully"
	lang["catDelSucc"] = "Category deleted successfully"
	lang["catall"] = "Category list"
	lang["CategoryModifiedSucc"] = "Category modified successfully"
	lang["showname"] = "Display name"
	lang["path"] = "Path"
	lang["EnterDisplayName"] = "Please enter a display name:"
	lang["EnterPath"] = "Please enter a path:"
	lang["ConfirmModifyCategory"] = "Confirm to modify the category?"
	lang["AreYouSureAddCategory"] = "Are you sure to add a new category?"
	lang["EnterCatIDDel"] = "Please enter the category ID to be deleted:"
	lang["moveCatId"] = "Please enter the category ID to which the article under the current category is moved:"
	lang["AreYouSureDelCategory"] = "Are you sure to delete the category?"
	lang["EnterCategoryIDModified"] = "Please enter the category ID to be modified:"
}