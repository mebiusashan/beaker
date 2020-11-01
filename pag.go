package beaker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/apcera/termtables"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/russross/blackfriday"
)

type PageDB struct {
	gorm.Model
	Title   string
	Context string
}

func (PageDB) TableName() string {
	return "pages"
}

type PagMDao struct {
	baseModel
}

func newPage(d *dao) *PagMDao {
	p := new(PagMDao)
	p.blogDao = d
	return p
}

func (p *PagMDao) FindAll() ([]PageDB, error) {
	var pages []PageDB
	err := p.blogDao.mysql.Order("id").Select("id, title").Find(&pages).Error
	return pages, err
}

func (p *PagMDao) FindByID(id uint) (*PageDB, error) {
	page := new(PageDB)
	err := p.blogDao.mysql.Where("id = ?", id).First(page).Error
	return page, err
}

func (p *PagMDao) Add(title string, context string) error {
	page := new(PageDB)
	page.Title = title
	page.Context = context
	err := p.blogDao.mysql.Create(&page).Error
	return err
}

func (p *PagMDao) Del(id uint) error {
	page := PageDB{}
	page.ID = id
	err := p.blogDao.mysql.Delete(&page).Error
	return err
}

func (ct *PagCtrl) Do(c *gin.Context) {
	idstr := c.Param("id")

	id, err := strconv.Atoi(idstr)
	if err != nil {
		ct.ctrl.ErrC.Do404(c)
		return
	}

	bodyStr, err := ct.ctrl.mvc.cache.GET(TAG_PAGE, idstr)

	if err == nil && bodyStr != "" {
		write200(c, bodyStr)
		return
	}

	page, err := ct.ctrl.mvc.model.PagDao.FindByID(uint(id))
	if err != nil {
		ct.ctrl.ErrC.Do500(c)
		return
	}
	bodyStr = string(blackfriday.Run([]byte(page.Context)))

	vars := ct.ctrl.mvc.view.GetVarMap()
	vars.Set("body", bodyStr)
	vars.Set("title", page.Title)

	bodyStr, err = ct.ctrl.mvc.view.Render(PAGE, vars)
	if err != nil {
		ct.ctrl.ErrC.Do500(c)
		return
	}

	ct.ctrl.mvc.cache.SETNX(TAG_PAGE, idstr, bodyStr, ct.ctrl.mvc.config.Redis.EXPIRE_TIME)
	write200(c, bodyStr)
}

func (ct *PagCtrl) Add(c *gin.Context) {
	var postData PageDB
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}
	err = ct.ctrl.mvc.model.PagDao.Add(postData.Title, postData.Context)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	ct.ctrl.mvc.cache.ClearAll()
	writeSucc(c, "单页添加成功", nil)
}

func (ct *PagCtrl) Del(c *gin.Context) {
	var postData PageDB
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}
	err = ct.ctrl.mvc.model.PagDao.Del(postData.ID)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	ct.ctrl.mvc.cache.ClearAll()
	writeSucc(c, "单页删除成功", nil)
}

func (ct *PagCtrl) List(c *gin.Context) {

	pags, err := ct.ctrl.mvc.model.PagDao.FindAll()
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	writeSucc(c, "", pags)
}

func (ct *PagCtrl) Down(c *gin.Context) {
	var postData PageDB
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}
	page, err := ct.ctrl.mvc.model.PagDao.FindByID(postData.ID)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	writeSucc(c, "", page)
}

func CMDPagAll() {
	resp, err := http.Post(HOST+"/admin/pag/list", "", strings.NewReader(""))
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

	table := termtables.CreateTable()
	table.AddHeaders("ID", "标题")
	for _, va := range jsonData.Data.([]interface{}) {
		v := va.(map[string]interface{})
		table.AddRow(uint(v["ID"].(float64)), v["Title"])
	}
	fmt.Println(table.Render())
}

func CMDPagDel() {
	fmt.Printf("请输入要删除的单页ID：")
	var delid uint
	fmt.Scanln(&delid)

	if delid == 0 {
		fmt.Println("ID错误")
		return
	}

	fmt.Printf("确认删除单页么？" +
		"\n分类ID：\"" + strconv.Itoa(int(delid)) + "\"" +
		"\n确认输入（y or n）:")
	var yes string
	fmt.Scanln(&yes)

	if yes != "y" && yes != "yes" {
		return
	}

	sendData := ArcDB{}
	sendData.ID = delid
	jsonByte, err := json.Marshal(sendData)
	if err != nil {
		//fmt.Println("login", err)
	}

	resp, err := http.Post(HOST+"/admin/pag/del", "", strings.NewReader(string(jsonByte)))
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

	fmt.Println("单页删除成功")
}

func CMDPagDown() {
	fmt.Printf("请输入要下载的单页ID：")
	var delid uint
	fmt.Scanln(&delid)

	if delid == 0 {
		fmt.Println("ID错误")
		return
	}

	fmt.Printf("确认下载单页么？" +
		"\n分类ID：\"" + strconv.Itoa(int(delid)) + "\"" +
		"\n确认输入（y or n）:")
	var yes string
	fmt.Scanln(&yes)

	if yes != "y" && yes != "yes" {
		return
	}

	sendData := ArcDB{}
	sendData.ID = delid
	jsonByte, err := json.Marshal(sendData)
	if err != nil {
		//fmt.Println("login", err)
	}

	resp, err := http.Post(HOST+"/admin/pag/down", "", strings.NewReader(string(jsonByte)))
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

	data := jsonData.Data.(map[string]interface{})
	path, _ := os.Getwd()
	path = path + "/" + strconv.Itoa(int(data["ID"].(float64))) + "_" + data["Title"].(string) + ".md"
	str := data["Context"].(string)
	ioutil.WriteFile(path, []byte(str), 0666)
	fmt.Println("文章下载成功，存储到：", path)
}

func CMDPagAdd() {

	if len(os.Args) < 2 {
		fmt.Println("缺少md文件")
		return
	}

	mdPath := os.Args[1]
	if mdPath == "" {
		fmt.Println("没有找到md文件")
		return
	}

	has, err := PathExists(mdPath)
	if err != nil || !has {
		fmt.Println("没有找到md文件")
		return
	}

	fmt.Printf("请输入单页标题：")

	reader := bufio.NewReader(os.Stdin)

	title, _, _ := reader.ReadLine()

	fmt.Printf("确认添加Tweet么？" +
		"\n单页标题：\"" + string(title) + "\"" +
		"\n内容文件：\"" + string(mdPath) + "\"" +
		"\n确认输入（y or n）:")
	var yes string
	fmt.Scanln(&yes)

	if yes != "y" && yes != "yes" {
		return
	}

	context, err := ioutil.ReadFile(mdPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	sendData := PageDB{Title: string(title), Context: string(context)}
	jsonByte, err := json.Marshal(sendData)
	if err != nil {
		//fmt.Println("login", err)
	}

	resp, err := http.Post(HOST+"/admin/pag/add", "", strings.NewReader(string(jsonByte)))
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

	fmt.Println("单页添加成功")
}
