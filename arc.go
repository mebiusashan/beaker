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

type ArcDB struct {
	gorm.Model
	Catid   uint
	Title   string
	Context string
}

func (ArcDB) TableName() string {
	return "arcs"
}

type ArcMDao struct {
	baseModel
}

func newArc(d *dao) *ArcMDao {
	a := new(ArcMDao)
	a.blogDao = d
	return a
}

func (a *ArcMDao) FindsByNum(num uint) ([]ArcDB, error) {
	var arcs []ArcDB
	err := a.blogDao.mysql.Order("id desc").Select("id, title, created_at").Limit(num).Find(&arcs).Error
	return arcs, err
}

func (a *ArcMDao) FindByID(id uint) (*ArcDB, error) {
	arc := new(ArcDB)
	err := a.blogDao.mysql.Where("id = ?", id).First(arc).Error
	return arc, err
}

func (a *ArcMDao) FindByCatID(catid uint) ([]ArcDB, error) {
	var arcs []ArcDB
	err := a.blogDao.mysql.Order("id desc").Select("id, title, created_at").Where("catid = ?", catid).Find(&arcs).Error
	return arcs, err
}

func (a *ArcMDao) FindAll() ([]ArcDB, error) {
	var arcs []ArcDB
	err := a.blogDao.mysql.Select("id, title, created_at, updated_at, catid").Find(&arcs).Error
	return arcs, err
}

func (a *ArcMDao) Add(catid uint, title string, context string) error {
	arc := new(ArcDB)
	arc.Catid = catid
	arc.Title = title
	arc.Context = context
	err := a.blogDao.mysql.Create(&arc).Error
	return err
}

func (a *ArcMDao) Del(id uint) error {
	arc := new(ArcDB)
	arc.ID = id
	err := a.blogDao.mysql.Delete(&arc).Error
	return err
}

func (a *ArcMDao) UpdateCat(catID uint, mvCatId uint) error {
	arc := new(ArcDB)
	err := a.blogDao.mysql.Model(&arc).Where("catid = ?", catID).Update("catid", mvCatId).Error
	return err
}

func (ct *ArcCtrl) Do(c *gin.Context) {
	idstr := c.Param("id")

	id, err := strconv.Atoi(idstr)
	if err != nil {
		ct.ctrl.ErrC.Do404(c)
		return
	}

	bodyStr, err := ct.ctrl.mvc.cache.GET(TAG_ARCHIVE, idstr)

	if err == nil && bodyStr != "" {
		write200(c, bodyStr)
		return
	}

	arcs, err := ct.ctrl.mvc.model.ArcDao.FindByID(uint(id))
	if err != nil {
		ct.ctrl.ErrC.Do500(c)
		return
	}
	bodyStr = string(blackfriday.Run([]byte(arcs.Context)))

	vars := ct.ctrl.mvc.view.GetVarMap()
	vars.Set("body", bodyStr)
	vars.Set("title", arcs.Title)

	bodyStr, err = ct.ctrl.mvc.view.Render(ARCHIVE, vars)
	if err != nil {
		ct.ctrl.ErrC.Do500(c)
		return
	}

	ct.ctrl.mvc.cache.SETNX(TAG_ARCHIVE, idstr, bodyStr, ct.ctrl.mvc.config.Redis.EXPIRE_TIME)
	write200(c, bodyStr)
}

func (ct *ArcCtrl) Add(c *gin.Context) {
	var postData ArcDB
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	err = ct.ctrl.mvc.model.ArcDao.Add(postData.Catid, postData.Title, postData.Context)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	ct.ctrl.mvc.cache.ClearAll()
	writeSucc(c, "文章添加成功", nil)
}

func (ct *ArcCtrl) Del(c *gin.Context) {
	var postData ArcDB
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	err = ct.ctrl.mvc.model.ArcDao.Del(postData.ID)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	ct.ctrl.mvc.cache.ClearAll()
	writeSucc(c, "文章删除成功", nil)
}

func (ct *ArcCtrl) All(c *gin.Context) {
	arcs, err := ct.ctrl.mvc.model.ArcDao.FindAll()
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	writeSucc(c, "全部文章", arcs)
}

func (ct *ArcCtrl) Down(c *gin.Context) {
	var postData ArcDB
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	arc, err := ct.ctrl.mvc.model.ArcDao.FindByID(postData.ID)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	writeSucc(c, "文章", arc)
}

func CMDArcAll() {
	resp, err := http.Post(HOST+"/admin/arc/list", "", strings.NewReader(""))
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

func CMDArcDel() {
	fmt.Printf("请输入要删除的文章ID：")
	var delid uint
	fmt.Scanln(&delid)

	if delid == 0 {
		fmt.Println("ID错误")
		return
	}

	fmt.Printf("确认删除文章么？" +
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

	resp, err := http.Post(HOST+"/admin/arc/del", "", strings.NewReader(string(jsonByte)))
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

	fmt.Println("文章删除成功")
}

func CMDArcDown() {
	fmt.Printf("请输入要下载的文章ID：")
	var delid uint
	fmt.Scanln(&delid)

	if delid == 0 {
		fmt.Println("ID错误")
		return
	}

	fmt.Printf("确认下载文章么？" +
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

	resp, err := http.Post(HOST+"/admin/arc/down", "", strings.NewReader(string(jsonByte)))
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

func CMDArcAdd() {

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

	fmt.Printf("请输入文章分类ID：")
	var delid uint
	fmt.Scanln(&delid)

	if delid == 0 {
		fmt.Println("ID错误")
		return
	}

	fmt.Printf("请输入单页标题：")

	reader := bufio.NewReader(os.Stdin)

	title, _, _ := reader.ReadLine()

	fmt.Printf("确认添加Tweet么？" +
		"\n分类ID：\"" + strconv.Itoa(int(delid)) + "\"" +
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

	sendData := ArcDB{Title: string(title), Context: string(context)}
	sendData.Catid = delid
	jsonByte, err := json.Marshal(sendData)
	if err != nil {
		//fmt.Println("login", err)
	}

	resp, err := http.Post(HOST+"/admin/arc/add", "", strings.NewReader(string(jsonByte)))
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

	fmt.Println("文章添加成功")
}
