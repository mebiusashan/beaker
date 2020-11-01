package beaker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/apcera/termtables"
	"github.com/gin-gonic/gin"
)

func (ct *OptCtrl) Info(c *gin.Context) {
	path := ct.ctrl.mvc.config.AuthInfo.ConfigPath
	config, err := NewWithPath(path, 0x1B)
	if err != nil {
		writeFail(c, err.Error())
		return
	}

	writeSucc(c, "website", config.Website)
}

func (ct *OptCtrl) ClearCache(c *gin.Context) {
	ct.ctrl.mvc.cache.ClearAll()
	writeSucc(c, "清除缓存成功", nil)
}

func CMDClearAllCache() {
	resp, err := http.Post(HOST+"/admin/clr/cache", "", strings.NewReader(""))
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

	fmt.Println("清除完成")
}

func CMDSeeOpts() {
	resp, err := http.Post(HOST+"/admin/opt", "", strings.NewReader(""))
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
	table.AddHeaders("项目", "值")
	for k, v := range jsonData.Data.(map[string]interface{}) {
		table.AddRow(k, v)
		//fmt.Println(k, v)
	}
	fmt.Println(table.Render())

}
