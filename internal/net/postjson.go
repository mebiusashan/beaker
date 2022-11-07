package net

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/mebiusashan/beaker/internal/cert"
	"github.com/mebiusashan/beaker/internal/common"
)

func PostJsonWithEncrypt(url string, refresh bool, key []byte, data interface{}) common.SuccMsgResp {
	var postData common.BaseReqMsg
	postData.Refresh = refresh
	postData.Data = data
	jsonByte, err := json.Marshal(postData)
	common.Assert(err)
	desT, err := cert.TripleDesEncrypt(jsonByte, key)
	return PostJson(url, strings.NewReader(cert.Base64Encode(desT)))
}

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
