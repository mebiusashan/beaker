package beaker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/apcera/termtables"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CatDB struct {
	gorm.Model
	Cname string
	Name  string
}

func (CatDB) TableName() string {
	return "cats"
}

type CatMDao struct {
	baseModel
}

func newCat(d *dao) *CatMDao {
	c := new(CatMDao)
	c.blogDao = d
	return c
}

func (c *CatMDao) FindAll() ([]CatDB, error) {
	var cats []CatDB
	err := c.blogDao.mysql.Order("id desc").Find(&cats).Error
	return cats, err
}

func (c *CatMDao) FindByCName(name string) (*CatDB, error) {
	cat := new(CatDB)
	err := c.blogDao.mysql.Order("id desc").Where("cname = ?", name).First(cat).Error
	return cat, err
}

func (c *CatMDao) Add(name string, cname string) error {
	cat := CatDB{Name: name, Cname: cname}
	err := c.blogDao.mysql.Create(&cat).Error
	return err
}

func (c *CatMDao) Del(id uint) error {
	cat := CatDB{}
	cat.ID = id
	err := c.blogDao.mysql.Delete(cat).Error
	return err
}

func (c *CatMDao) Update(id uint, name string, cname string) error {
	data := make(map[string]interface{})
	if name != "" {
		data["name"] = name
	}
	if cname != "" {
		data["cname"] = cname
	}
	cat := CatDB{}
	cat.ID = id
	err := c.blogDao.mysql.Model(&cat).Updates(data).Error
	return err
}

func (ct *CatCtrl) Do(c *gin.Context) {
	cname := c.Param("name")
	bodyStr, err := ct.ctrl.mvc.cache.GET(TAG_CAT, cname)

	if err == nil && bodyStr != "" {
		write200(c, bodyStr)
		return
	}

	pages, err := ct.ctrl.mvc.model.PagDao.FindAll()
	if err != nil {
		fmt.Println(1, err)
		ct.ctrl.ErrC.Do500(c)
		return
	}

	cats, err := ct.ctrl.mvc.model.CatDao.FindAll()
	if err != nil {
		ct.ctrl.ErrC.Do500(c)
		return
	}

	cat, err := ct.ctrl.mvc.model.CatDao.FindByCName(cname)
	if err != nil {
		fmt.Println(2, cname, err)
		ct.ctrl.ErrC.Do500(c)
		return
	}

	fmt.Println(cat.ID)
	arcs, err := ct.ctrl.mvc.model.ArcDao.FindByCatID(cat.ID)
	if err != nil {
		fmt.Println(3, err)
		ct.ctrl.ErrC.Do500(c)
		return
	}

	vars := ct.ctrl.mvc.view.GetVarMap()
	vars.Set("title", cat.Name)
	vars.Set("cats", cats)
	vars.Set("pages", pages)
	vars.Set("arcs", arcs)

	bodyStr, err = ct.ctrl.mvc.view.Render(CAT, vars)
	if err != nil {
		fmt.Println(4, err)
		ct.ctrl.ErrC.Do500(c)
		return
	}

	ct.ctrl.mvc.cache.SETNX(TAG_CAT, cname, bodyStr, ct.ctrl.mvc.config.Redis.EXPIRE_TIME)
	write200(c, bodyStr)
}

func (ct *CatCtrl) Add(c *gin.Context) {
	var postData CatDB
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}
	if postData.Name == "" || postData.Cname == "" {
		writeFail(c, "不允许出现空值")
		return
	}

	err = ct.ctrl.mvc.model.CatDao.Add(postData.Name, postData.Cname)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	ct.ctrl.mvc.cache.ClearAll()
	writeSucc(c, "分类添加成功", nil)
}

type CatDBDel struct {
	CatDB
	MvID uint
}

func (ct *CatCtrl) Del(c *gin.Context) {
	var postData CatDBDel
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	err = ct.ctrl.mvc.model.ArcDao.UpdateCat(postData.ID, postData.MvID)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	err = ct.ctrl.mvc.model.CatDao.Del(postData.ID)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	ct.ctrl.mvc.cache.ClearAll()
	writeSucc(c, "分类删除成功", nil)
}

func (ct *CatCtrl) All(c *gin.Context) {
	list, err := ct.ctrl.mvc.model.CatDao.FindAll()
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	//ct.ctrl.mvcc.cache.ClearAll()
	writeSucc(c, "所有分了", list)
}

func (ct *CatCtrl) Update(c *gin.Context) {
	var postData CatDB
	err := c.BindJSON(&postData)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	err = ct.ctrl.mvc.model.CatDao.Update(postData.ID, postData.Name, postData.Cname)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	ct.ctrl.mvc.cache.ClearAll()
	writeSucc(c, "分类修改成功", nil)
}

func CMDCatAll() {
	resp, err := http.Post(HOST+"/admin/cat/list", "", strings.NewReader(""))
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
	table.AddHeaders("ID", "显示名", "路径名", "创建时间")
	for _, va := range jsonData.Data.([]interface{}) {
		v := va.(map[string]interface{})
		table.AddRow(uint(v["ID"].(float64)), v["Name"], v["Cname"], v["CreatedAt"])
	}
	fmt.Println(table.Render())
}

func CMDCatAdd() {

	fmt.Printf("请输入显示名：")
	var name string
	fmt.Scanln(&name)

	fmt.Printf("请输入路径名：")
	var cname string
	fmt.Scanln(&cname)

	fmt.Printf("确认添加新分类么？" +
		"\n显示名：\"" + name + "\"" +
		"\n路径名:\"" + cname + "\"" +
		"\n确认输入（y or n）:")
	var yes string
	fmt.Scanln(&yes)

	if yes != "y" && yes != "yes" {
		return
	}

	sendData := CatDB{Cname: cname, Name: name}
	jsonByte, err := json.Marshal(sendData)
	if err != nil {
		//fmt.Println("login", err)
	}

	resp, err := http.Post(HOST+"/admin/cat/add", "", strings.NewReader(string(jsonByte)))
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

	fmt.Println("分类添加成功")
}

func CMDCatDel() {
	fmt.Printf("请输入要删除的分类ID：")
	var delid uint
	fmt.Scanln(&delid)

	if delid == 0 {
		fmt.Println("ID错误")
		return
	}

	fmt.Printf("请输入当前分类下文章移动到的分类ID：")
	var mvid uint
	fmt.Scanln(&mvid)

	if mvid == 0 {
		fmt.Println("ID错误")
		return
	}

	fmt.Printf("确认删除分类么？" +
		"\n分类ID：\"" + strconv.Itoa(int(delid)) + "\"" +
		"\n确认输入（y or n）:")
	var yes string
	fmt.Scanln(&yes)

	if yes != "y" && yes != "yes" {
		return
	}

	sendData := CatDBDel{}
	sendData.ID = delid
	sendData.MvID = mvid
	jsonByte, err := json.Marshal(sendData)
	if err != nil {
		//fmt.Println("login", err)
	}

	resp, err := http.Post(HOST+"/admin/cat/del", "", strings.NewReader(string(jsonByte)))
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

	fmt.Println("分类删除成功")
}

func CMDCatEdit() {
	fmt.Printf("请输入要修改的分类ID：")
	var delid uint
	fmt.Scanln(&delid)

	if delid == 0 {
		fmt.Println("ID错误")
		return
	}

	fmt.Printf("请输入显示名：")
	var name string
	fmt.Scanln(&name)

	fmt.Printf("请输入路径名：")
	var cname string
	fmt.Scanln(&cname)

	fmt.Printf("确认修改分类么？" +
		"\n分类ID：\"" + strconv.Itoa(int(delid)) + "\"" +
		"\n显示名：\"" + name + "\"" +
		"\n路径名:\"" + cname + "\"" +
		"\n确认输入（y or n）:")
	var yes string
	fmt.Scanln(&yes)

	if yes != "y" && yes != "yes" {
		return
	}

	sendData := CatDB{Cname: cname, Name: name}
	sendData.ID = delid
	jsonByte, err := json.Marshal(sendData)
	if err != nil {
		//fmt.Println("login", err)
	}

	resp, err := http.Post(HOST+"/admin/cat/change", "", strings.NewReader(string(jsonByte)))
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

	fmt.Println("分类修改成功")
}
