package net

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mebiusashan/beaker/common"
)

func PostJson(url string, body io.Reader) common.SuccMsgResp {
	resp, err := http.Post(url, "", body)
	common.Assert(err)

	Body, err := ioutil.ReadAll(resp.Body)
	common.Assert(err)

	var jsonData common.SuccMsgResp
	err = json.Unmarshal(Body, &jsonData)
	if err != nil {
		fmt.Println(string(Body))
		os.Exit(0)
	}

	if jsonData.Code != common.SUCC {
		common.Err(jsonData.Msg)
	}
	return jsonData
}
