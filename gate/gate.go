package gate

import (
    "encoding/json"
    "encoding/xml"
    "net/http"
)

type Gate struct {
    W http.ResponseWriter
}

type RegRoute struct {
    Auth   string `json:"auth"`
    Desc   string `json:"desc"`
    Method string `json:"method"`
    Uri    string `json:"uri"`
    Host   string `json:"host"`
}

type ApiData struct {
    Module    string                 `json:"module"`
    Uri       string                 `json:"uri"`
    UserToken string                 `json:"user-token"`
    UserInfo  map[string]interface{} `json:"user_info"`
    Param     map[string]interface{} `json:"param"`
    RequestId string                 `json:"request_id"`
    Type      string                 `json:"type"`
    IP        string                 `json:"IP"`
}

// 跨域
func (gate Gate) Cross() {
    gate.W.Header().Set("Access-Control-Allow-Origin", "*")
    gate.W.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    gate.W.Header().Set("Access-Control-Allow-Headers", "Action, ApiAuth, User-Token, Module, X-PINGOTHER, Content-Type, Content-Disposition")
    gate.W.Header().Set("Content-Type", "application/json; charset=utf-8")
}

// 返回json
func (gate Gate) Json(statusCode int, data interface{}) {
    gate.Cross()
    gate.W.WriteHeader(statusCode)
    str, _ := json.Marshal(data)
    _, _ = gate.W.Write(str)
}

// 返回xml
func (gate Gate) Xml(statusCode int, data interface{}) {
    gate.Cross()
    gate.W.Header().Set("Content-Type", "application/xml; charset=utf-8")
    gate.W.WriteHeader(statusCode)
    str, _ := xml.Marshal(data)
    _, _ = gate.W.Write(str)
}
