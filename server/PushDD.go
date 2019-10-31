package server

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strings"
)

func PushToDingDing(accessToken, content string) {
    var (
        url     string
        data    map[string]interface{}
        data1   map[string]interface{}
        str     map[string]interface{}
        resMap  map[string]interface{}
        byteStr []byte
        resp    *http.Response
    )
    url = "https://oapi.dingtalk.com/robot/send?access_token=" + accessToken
    data = map[string]interface{}{}
    data1 = map[string]interface{}{}
    str = map[string]interface{}{}
    str["msgtype"] = "text"
    data["content"] = content
    data1["isAtAll"] = false
    str["text"] = data
    str["at"] = data1
    byteStr, _ = json.Marshal(&str)
    resp, _ = http.Post(url, "application/json", strings.NewReader(string(byteStr)))
    res, _ := ioutil.ReadAll(resp.Body)
    err := json.Unmarshal(res, &resMap)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
}
