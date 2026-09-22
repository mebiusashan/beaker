package net

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mebiusashan/beaker/internal/cert"
	"github.com/mebiusashan/beaker/internal/common"
)

var session struct {
	sync.RWMutex
	token string
}

func SetSessionToken(token string) {
	session.Lock()
	session.token = token
	session.Unlock()
}

func SessionToken() string {
	session.RLock()
	defer session.RUnlock()
	return session.token
}

func PostJsonWithEncrypt(url string, refresh bool, key []byte, data interface{}) common.SuccMsgResp {
	var postData common.BaseReqMsg
	postData.Refresh = refresh
	postData.Data = data
	jsonByte, err := json.Marshal(postData)
	if err != nil {
		common.ErrCode(common.ErrorCodeInvalidRequest, err.Error())
	}
	desT, err := cert.TripleDesEncrypt(jsonByte, key)
	if err != nil {
		common.ErrCode(common.ErrorCodeInvalidRequest, err.Error())
	}
	return PostJson(url, strings.NewReader(cert.Base64Encode(desT)))
}

func PostJson(url string, body io.Reader) common.SuccMsgResp {
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		common.ErrCode(common.ErrorCodeNetwork, err.Error())
	}
	request.Header.Set("Content-Type", "application/octet-stream")
	request.Header.Set("Accept", "application/json")
	if token := SessionToken(); token != "" {
		request.Header.Set("X-Beaker-Session", token)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(request)
	if err != nil {
		common.ErrCode(common.ErrorCodeNetwork, err.Error())
	}
	defer resp.Body.Close()

	const maxResponseSize = 10 << 20
	limited := io.LimitReader(resp.Body, maxResponseSize)
	Body, err := io.ReadAll(limited)
	if err != nil {
		common.ErrCode(common.ErrorCodeNetwork, err.Error())
	}

	var jsonData common.SuccMsgResp
	err = json.Unmarshal(Body, &jsonData)
	if err != nil {
		common.ErrCode(common.ErrorCodeInvalidResponse, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(Body)))
	}

	if jsonData.Code != common.SUCC {
		if jsonData.ErrorCode == "" {
			jsonData.ErrorCode = common.ErrorCodeHTTP
		}
		common.ErrCodeWithRequest(jsonData.ErrorCode, jsonData.Msg, jsonData.RequestID)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		code := jsonData.ErrorCode
		if code == "" {
			code = common.ErrorCodeHTTP
		}
		common.ErrCodeWithRequest(code, jsonData.Msg, jsonData.RequestID)
	}
	return jsonData
}
