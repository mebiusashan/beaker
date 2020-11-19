package net

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/mebiusashan/beaker/common"
)

func PostJson(url string, body io.Reader) common.SuccMsg {
	resp, err := http.Post(url, "", body)
	common.Assert(err)

	Body, err := ioutil.ReadAll(resp.Body)
	common.Assert(err)

	var jsonData common.SuccMsg
	err = json.Unmarshal(Body, &jsonData)
	common.Assert(err)

	if jsonData.Code != common.SUCC {
		common.Err(jsonData.Msg)
	}
	return jsonData
}
