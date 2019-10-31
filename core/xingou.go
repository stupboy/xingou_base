package core

import (
    "bufio"
    "encoding/json"
    "fmt"
    "github.com/garyburd/redigo/redis"
    "github.com/xingou_base/cache"
    "github.com/xingou_base/gate"
    "github.com/xingou_base/utli"
    "github.com/xingou_base/validate"
    "io"
    "log"
    "net"
    "os"
    "strconv"
    "strings"
    "time"
)

// 路由存储对象
var FuncMap map[string]interface{}
var ValidateMap map[string]interface{}
var Route XinGou
var RouteInfo map[string]interface{}
var MiddleFunc []func(route XinGou, data gate.ApiData) (XinGou)

type XinGou struct {
    Writer    *net.UDPConn
    Addr      net.UDPAddr
    Data      map[string]interface{}
    UserToken string
    UserInfo  map[string]interface{}
    RequestId string
    ApiDoc    map[string]interface{}
    Url       string
    IP        string
}

// 初始化参数
func Init() {
    init := make(map[string]interface{})
    init["version"] = "v1.0.1"
    init["auth"] = "king"
    initInfo := make(map[string]interface{})
    FuncMap = init
    RouteInfo = initInfo
    MiddleFunc = make([]func(gou XinGou, data gate.ApiData) (XinGou), 0)
}

// 返写数据
func (c *XinGou) Return(data []byte) {
    // 验证输出字段
    if len(data) >= 128 {
        r := cache.RedisPool.Get()
        defer r.Close()
        _, _ = r.Do("set", c.RequestId, data)
        _, _ = r.Do("expire", c.RequestId, 5)
        back := make(map[string]string)
        back["type"] = "uuid"
        back["response_id"] = c.RequestId
        backJson, _ := json.Marshal(back)
        _, _ = c.Writer.WriteToUDP(backJson, &c.Addr)
    } else {
        _, _ = c.Writer.WriteToUDP(data, &c.Addr)
    }
}

func (c *XinGou) GetInt(key string, defaultValue ...int) interface{} {
    value, ok := c.Data[key]
    if !ok {
        if defaultValue != nil {
            return defaultValue[0]
        }
        return nil
    }
    ret, _ := strconv.Atoi(value.(string))
    return ret
}

func (c *XinGou) GetStringSlice(key string, defaultValue ...int) interface{} {
    value, ok := c.Data[key]
    if !ok {
        if defaultValue != nil {
            return defaultValue[0]
        }
        return nil
    }
    return value.([]interface{})
}

func (c *XinGou) GetFloat64(key string, defaultValue ...float64) interface{} {
    value, ok := c.Data[key]
    if !ok {
        if defaultValue != nil {
            return defaultValue[0]
        }
        return nil
    }
    return value
}

func (c *XinGou) GetString(key string, defaultValue ...string) interface{} {
    value, ok := c.Data[key]
    if !ok {
        if defaultValue != nil {
            return defaultValue[0]
        }
        return nil
    }
    ret, _ := value.(string)
    return ret
}

func (c *XinGou) ReturnSuccess(ret ...interface{}) {
    data := make(map[string]interface{})
    data["code"] = utli.ReturnSuccess
    data["msg"] = utli.ErrCodeMap[utli.ReturnSuccess]
    if len(ret) > 1 {
        data["data"] = ret
    }
    if len(ret) == 1 {
        data["data"] = ret[0]
    }
    //返回参数校验
    if _, ok := c.ApiDoc[c.Url].(map[string]interface{})["return"]; ok {
        // 有填写返回注释才校验返回值
        returnRule := c.ApiDoc[c.Url].(map[string]interface{})["return"]
        validate.ReturnValidate(data, returnRule.(map[string]interface{}))
    }
    jsonData, _ := json.Marshal(data)
    // 验证输出字段
    c.Return(jsonData)
}

func (c *XinGou) ReturnError(errCode int, err error) {
    data := make(map[string]interface{})
    data["code"] = utli.ReturnError
    data["error_code"] = errCode
    data["msg"] = utli.ErrCodeMap[errCode]
    data["error"] = err
    retData, _ := json.Marshal(data)
    c.Return(retData)
}

func InitValidate() {
    baseValidate := make(map[string]interface{})
    ValidateMap = baseValidate
    ValidateMap["len"] = validate.LenValidate
    ValidateMap["max"] = validate.MaxValidate
    ValidateMap["min"] = validate.MinValidate
    ValidateMap["range"] = validate.RangValidate
    ValidateMap["mobile"] = validate.MobileValidate
    ValidateMap["in"] = validate.InValidate
}

func AddRoute(name string, method func(c *XinGou), desc ...string) {
    // 注册路由
    FuncMap[name] = method
    // 注册信息
    temp := make(map[string]string)
    temp["uri"] = name
    // 路由默认字段注册
    for key, val := range desc {
        strMap := map[int]string{0: "desc", 1: "method", 2: "auth"}
        defaultMap := map[int]string{0: "服务描述", 1: "ANY", 2: "NONE"}
        if val != "" {
            temp[strMap[key]] = desc[key]
            if Route.ApiDoc[name] != nil {
                Route.ApiDoc[name].(map[string]interface{})[strMap[key]] = desc[key]
            }
        } else {
            temp[strMap[key]] = defaultMap[key]
        }
    }
    RouteInfo[name] = temp
}


func NoRoute(method func(c *XinGou)) {
    log.Println("添加MISS路由")
    FuncMap["NoRoute"] = method
}

func HandleConn(listener *net.UDPConn, Addr net.UDPAddr, data []byte) {
    var (
        jsonData gate.ApiData
        r        redis.Conn
        backData map[string]interface{}
        ok       bool
        err      error
        dd       []byte
    )
    RouteTemp := Route
    RouteTemp.Writer = listener
    RouteTemp.Addr = Addr
    apiTime := time.Now().UnixNano()
    // 全局异常处理
    defer func() {
        if r := recover(); r != nil {
            var errMsg string
            switch val := r.(type) {
            case string:
                errMsg = val
            case error:
                errMsg = val.Error()
            default:
                errMsg = "未知错误类型"
            }
            backData = make(map[string]interface{})
            backData["code"] = utli.ReturnError
            backData["msg"] = errMsg
            backJson, _ := json.Marshal(backData)
            _, _ = RouteTemp.Writer.WriteToUDP(backJson, &RouteTemp.Addr)
        }
        if jsonData.Type != "ping" {
            expire := (time.Now().UnixNano() - apiTime) / 1e6
            log.Println("完成:", RouteTemp.Url+" ", expire, "ms")
        }
    }()
    // 请求数据
    _ = json.Unmarshal(data, &jsonData) // 替换
    if jsonData.Type == "uuid" {
        r = cache.RedisPool.Get()
        defer r.Close()
        dd, err = redis.Bytes(r.Do("get", jsonData.RequestId))
        if err != nil {
            panic("数据传输错误")
        }
        jsonData = gate.ApiData{}
        _ = json.Unmarshal(dd, &jsonData)
    }
    if jsonData.Type == "ping" {
        backData = make(map[string]interface{})
        backData["type"] = "pong"
        backJson, _ := json.Marshal(backData)
        _, _ = RouteTemp.Writer.WriteToUDP(backJson, &RouteTemp.Addr)
        // 记录心跳时间
        utli.PingMap.Store("ping_time", time.Now().Unix())
        return
    }
    // 参数写入
    RouteTemp.RequestId = jsonData.RequestId
    if jsonData.IP != "" {
        RouteTemp.IP = strings.Split(jsonData.IP, ":")[0]
    }
    // 参数写入
    if jsonData.Param != nil {
        RouteTemp.Data = utli.Float64ToInt(jsonData.Param, RouteTemp.ApiDoc[jsonData.Module+"/"+jsonData.Uri])
    }
    // token写入
    if jsonData.UserToken != "" {
        RouteTemp.UserToken = jsonData.UserToken
    }
    url := jsonData.Module + "/" + jsonData.Uri
    RouteTemp.Url = url
    log.Println("请求:", RouteTemp.Url, "请求Id:", RouteTemp.RequestId)
    // 参数验证
    if _, ok = RouteTemp.ApiDoc[url]; !ok {
        panic("接口注释未写")
    }
    log.Println(len(MiddleFunc))
    for _, f := range MiddleFunc {
        log.Println("test")
        RouteTemp = f(RouteTemp, jsonData)
    }
    RouteTemp.Data = validate.RequestValidate(RouteTemp.Data, RouteTemp.ApiDoc[url].(map[string]interface{}), ValidateMap)
    // 路由不存在的情况
    if _, ok = FuncMap[url]; !ok {
        if _, ok = FuncMap["NoRoute"]; ok {
            FuncMap["NoRoute"].(func(c *XinGou))(&RouteTemp)
        } else {
            panic("路由不存在")
        }
    } else {
        FuncMap[url].(func(c *XinGou))(&RouteTemp)
    }
}

func RegServer(RegHost string, ServerHost string) {
    var filename = "allow.txt"
    doc, err := os.Open(filename)
    allowMap := make(map[string]int)
    if err == nil {
        rd1 := bufio.NewReader(doc)
        for {
            line, err := rd1.ReadString('\n') //以'\n'为结束符读入一行
            line = strings.TrimSpace(line)
            if len(line) > 3 {
                allowMap[line] = 1
            }
            if err != nil || io.EOF == err {
                break
            }
        }
    }

    serverAddr := RegHost
    conn, err := net.Dial("udp", serverAddr)
    if err != nil {
        panic(err)
    }
    for key, val := range RouteInfo {
        if key == "version" {
            continue
        }
        if key == "auth" {
            continue
        }
        if key == "NoRoute" {
            continue
        }
        if len(allowMap) > 0 {
            if _, ok := allowMap[key]; !ok {
                continue
            }
        }
        log.Println("注册服务:", key)
        val.(map[string]string)["host"] = ServerHost
        toWrite, _ := json.Marshal(val)
        _, err := conn.Write([]byte(toWrite))
        if err != nil {
            panic(err)
        }
        msg := make([]byte, 100)
        err = conn.SetReadDeadline(time.Now().Add(time.Duration(5) * time.Second))
        if err != nil {
            log.Println(err)
            continue
        }
        n, err := conn.Read(msg)
        if err != nil {
            log.Println(err)
            panic(err)
        }
        if string(msg[0:n]) == "ok" {
            log.Println("注册成功")
        } else {
            log.Println("注册失败")
        }
    }
}

var udpLimit chan bool

func Run(HostParam ...string) {
    log.Println("监听:", HostParam[0])
    log.Println("Start Server Success!")
    addr, err := net.ResolveUDPAddr("udp", HostParam[0])
    udpLimit = make(chan bool, 1000)
    if err != nil {
        fmt.Print(err)
        return
    }
    listener, err := net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Print(err)
        return
    }
    defer listener.Close()
    for {
        udpLimit <- true
        go UdpServer(listener)
    }
}

func UdpServer(listener *net.UDPConn) {
    buf := make([]byte, 2048)
    n, ctlAddr, err := listener.ReadFromUDP(buf)
    if err != nil {
        fmt.Print(err)
        return
    }
    HandleConn(listener, *ctlAddr, buf[0:n])
    <-udpLimit
}
