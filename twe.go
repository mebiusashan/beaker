package beaker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/xlab/tablewriter"
)

type TweetDB struct {
	gorm.Model
	Context string
}

func (TweetDB) TableName() string {
	return "tweets"
}

type TweMDao struct {
	baseModel
}

func newTweet(d *dao) *TweMDao {
	t := new(TweMDao)
	t.blogDao = d
	return t
}

func (t *TweMDao) FindByNum(page uint, num uint) ([]TweetDB, error) {
	var tw []TweetDB
	err := t.blogDao.mysql.Order("id desc").Offset((page - 1) * num).Limit(num).Find(&tw).Error
	return tw, err
}

func (t *TweMDao) Count() uint {
	var count uint = 0
	t.blogDao.mysql.Table("tweets").Count(&count)
	return count
}

func (t *TweMDao) Add(context string) error {
	tw := TweetDB{Context: context}
	err := t.blogDao.mysql.Create(&tw).Error
	return err
}

func (t *TweMDao) Del(id uint) error {
	tw := TweetDB{}
	tw.ID = id
	err := t.blogDao.mysql.Delete(tw).Error
	return err
}

func (ct *TweCtrl) Do(c *gin.Context) {
	page := ct.convParam(c.Param("page"))
	numOfOnePage := ct.ctrl.mvc.config.Website.TWEET_NUM_ONE_PAGE

	twCount := ct.ctrl.mvc.model.TweDao.Count()
	twCount = uint(math.Ceil(float64(twCount) / float64(numOfOnePage)))
	if page > twCount {
		page = twCount
	}

	twpages := ct.createPageNum(twCount)

	tws, err := ct.ctrl.mvc.model.TweDao.FindByNum(page, numOfOnePage)
	if err != nil {
		ct.ctrl.ErrC.Do500(c)
		return
	}

	pages, err := ct.ctrl.mvc.model.PagDao.FindAll()
	if err != nil {
		ct.ctrl.ErrC.Do500(c)
		return
	}

	cats, err := ct.ctrl.mvc.model.CatDao.FindAll()
	if err != nil {
		ct.ctrl.ErrC.Do500(c)
		return
	}

	vars := ct.ctrl.mvc.view.GetVarMap()
	vars.Set("title", "Home")
	vars.Set("cats", cats)
	vars.Set("pages", pages)
	vars.Set("tws", tws)
	vars.Set("twpages", twpages)
	vars.Set("twcurpage", page)

	bodyStr, err := ct.ctrl.mvc.view.Render(TWEET, vars)
	if err != nil {
		ct.ctrl.ErrC.Do500(c)
		return
	}
	write200(c, bodyStr)
}

func (ct *TweCtrl) convParam(param string) uint {
	if param == "/" {
		return 1
	}
	str := param[1:]
	num, err := strconv.Atoi(str)
	if err != nil {
		return 1
	}
	return uint(num)
}

func (ct *TweCtrl) createPageNum(count uint) []uint {
	twpages := make([]uint, count)
	var i uint = 1
	for ; i <= count; i++ {
		twpages[i-1] = i
	}
	return twpages
}

func (ct *TweCtrl) Add(c *gin.Context) {
	var postData TweetDB
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}
	err = ct.ctrl.mvc.model.TweDao.Add(postData.Context)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	ct.ctrl.mvc.cache.ClearAll()
	writeSucc(c, GetLanguage("tweetAddSucc"), nil)
}

func (ct *TweCtrl) Del(c *gin.Context) {
	var postData TweetDB
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}
	err = ct.ctrl.mvc.model.TweDao.Del(postData.ID)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	ct.ctrl.mvc.cache.ClearAll()
	writeSucc(c, GetLanguage("tweetDelSucc"), nil)
}

type TweList struct {
	BaseMsg
	CurPage   uint
	TotlePage uint
	TweNum    uint
	List      []TweetDB
}

func (ct *TweCtrl) List(c *gin.Context) {
	var postData TweList
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	var page uint = 1
	if postData.CurPage >= 1 {
		page = postData.CurPage
	}

	twNums := ct.ctrl.mvc.model.TweDao.Count()
	postData.TweNum = twNums
	twNums = uint(math.Ceil(float64(twNums) / float64(10)))
	if page > twNums {
		page = twNums
	}

	tws, err := ct.ctrl.mvc.model.TweDao.FindByNum(page, 10)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	postData.Code = 0
	postData.CurPage = page
	postData.TotlePage = twNums
	postData.List = tws

	writeSucc(c, GetLanguage("tweetDelSucc"), postData)
}

func CMDTweAll() {

	postData := TweList{CurPage: 0}
	jsonByte, err := json.Marshal(postData)
	if err != nil {
		//fmt.Println("login", err)
	}

	resp, err := http.Post(HOST+"/admin/twe/list", "", strings.NewReader(string(jsonByte)))
	if err != nil {
		//fmt.Println("ping", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println("ping", err)
	}

	var jsonData SuccMsg
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		//fmt.Println("ping", err)
	}

	if jsonData.Code != SUCC {
		fmt.Println(jsonData.Msg)
		return
	}

	//fmt.Println(jsonData.Data)
	dd := jsonData.Data.(map[string]interface{})

	tablewriter.EnableUTF8()
	table := tablewriter.CreateTable()
	table.SetModeTerminal()
	table.AddHeaders("ID", GetLanguage("Content"), GetLanguage("CreateTime"))
	for _, v := range dd["List"].([]interface{}) {
		va := v.(map[string]interface{})
		table.AddRow(uint(va["ID"].(float64)), va["Context"], va["CreatedAt"])
	}
	fmt.Println(table.Render())
	fmt.Println(dd["TotlePage"], " pages,", dd["TweNum"], "tweets, current ", dd["CurPage"], "page")
}

func CMDTweDel() {
	fmt.Printf(GetLanguage("EnterTweetIDToDeleted"))
	var delid uint
	fmt.Scanln(&delid)

	if delid == 0 {
		fmt.Println(GetLanguage("IDError"))
		return
	}

	fmt.Printf(GetLanguage("AreYouSureDelTweet") +
		"\n"+GetLanguage("CategoryID")+":\"" + strconv.Itoa(int(delid)) + "\"" +
		"\n"+GetLanguage("ConfirmInput"))
	var yes string
	fmt.Scanln(&yes)

	if yes != "y" && yes != "yes" {
		return
	}

	sendData := TweetDB{}
	sendData.ID = delid
	jsonByte, err := json.Marshal(sendData)
	if err != nil {
		//fmt.Println("login", err)
	}

	resp, err := http.Post(HOST+"/admin/twe/del", "", strings.NewReader(string(jsonByte)))
	if err != nil {
		//fmt.Println("ping", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println("ping", err)
	}

	var jsonData SuccMsg
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		//fmt.Println("ping", err)
	}

	if jsonData.Code != SUCC {
		fmt.Println(jsonData.Msg)
		return
	}

	fmt.Println(GetLanguage("tweetDelSucc"))
}

func CMDTweAdd() {
	fmt.Printf(GetLanguage("EnterTweetContent"))

	reader := bufio.NewReader(os.Stdin)

	cname, _, _ := reader.ReadLine()

	//var cname string
	//fmt.Scanln(&cname)

	fmt.Printf(GetLanguage("AreYouSureAddTweet") +
		"\n"+GetLanguage("Content")+":\"" + string(cname) + "\"" +
		"\n"+GetLanguage("ConfirmInput"))
	var yes string
	fmt.Scanln(&yes)

	if yes != "y" && yes != "yes" {
		return
	}

	sendData := TweetDB{Context: string(cname)}
	jsonByte, err := json.Marshal(sendData)
	if err != nil {
		//fmt.Println("login", err)
	}

	resp, err := http.Post(HOST+"/admin/twe/add", "", strings.NewReader(string(jsonByte)))
	if err != nil {
		//fmt.Println("ping", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println("ping", err)
	}

	var jsonData SuccMsg
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		//fmt.Println("ping", err)
	}

	if jsonData.Code != SUCC {
		fmt.Println(jsonData.Msg)
		return
	}

	fmt.Println(GetLanguage("tweetAddSucc"))
}
