package xingou_base

import (
    "github.com/garyburd/redigo/redis"
    "github.com/stupboy/xingoudoc"
    "github.com/xingou_base/cache"
    "github.com/xingou_base/core"
    "github.com/xingou_base/gate"
    "github.com/xingou_base/utli"
    "log"
    "reflect"
)

type XinGouApp struct {
    Config    string      `json:"config"`
    RedisHost string      `json:"redis_host"`
    Dir       []string    `json:"dir"`
    C1        interface{} `json:"c_1"`
    Env       string
    RegHost   string
    SerHost   string
    Api       []interface{}
    Redis     *redis.Pool
    MiddleFunc []func(route core.XinGou, data gate.ApiData) (core.XinGou)
}

func (x *XinGouApp) AddMiddle(f func(gou core.XinGou, data gate.ApiData) (core.XinGou)) {
    if len(x.MiddleFunc) == 0 {
        x.MiddleFunc = []func(route core.XinGou, data gate.ApiData) (core.XinGou){f}
    } else {
        x.MiddleFunc = append(x.MiddleFunc, f)
    }
}

func (x *XinGouApp) GetApiDoc(urlDir ...string) map[string]interface{} {
    var err error
    //config := SystemConfig
    if x.Env == "pro" {
        return utli.GetApiDocJson()
    }
    if len(urlDir) == 0 {
        urlDir = []string{"controller/"}
    }
    var Doc xingoudoc.NoteDoc
    Doc.JsonName = "doc.json"
    for _, v := range urlDir {
        err = Doc.GetApiDoc(v + "/")
    }
    err = Doc.MapToJson()
    if err != nil {
        panic("error")
    }
    return Doc.Doc
}

func (x *XinGouApp) AddDir(dir string) {
    if len(x.Dir) == 0 {
        x.Dir = []string{dir}
    } else {
        x.Dir = append(x.Dir, dir)
    }
}

func (x *XinGouApp) AddApi(api interface{}) {
    if len(x.Api) == 0 {
        x.Api = []interface{}{api}
    } else {
        x.Api = append(x.Api, api)
    }
}

func (x *XinGouApp) AddValidate(key string, f func(val string, argusStr ...string) (bool)) {
    core.ValidateMap[key] = f
}

func (x *XinGouApp) AddConfig(file string) {
    x.Config = file
}

func (x *XinGouApp) InitRedis(host string, port string, auth string) {
    cache.InitRedis(host, port, auth)
    x.Redis = cache.RedisPool
}

func (x *XinGouApp) Server(regHost string, SerHost string) {
    x.RegHost = regHost
    x.SerHost = SerHost
}

func (x *XinGouApp) Run() {
    // 初始化
    core.Init()
    core.MiddleFunc = x.MiddleFunc
    core.Route.ApiDoc = x.GetApiDoc(x.Dir...)
    for _, api := range x.Api {
        data := reflect.ValueOf(api)
        for k, v := range core.Route.ApiDoc {
            route := v.(map[string]interface{})
            funcName := route["func"].(string)
            if data.MethodByName(funcName).IsValid() {
                log.Println(k, route["func"], "自动注解")
                core.AddRoute(k, func(c *core.XinGou) {
                    paramSlice := make([]reflect.Value, 1)
                    paramSlice[0] = reflect.ValueOf(c)
                    data.MethodByName(funcName).Call(paramSlice)
                }, route["title"].(string), route["method"].(string), route["auth"].(string))
            }
        }
    }
    // 初始化验证规则
    core.InitValidate()
    // 添加MISS路由
    // core.NoRoute(api.MissAction)
    // 监听UDP端口
    core.RegServer(x.RegHost, x.SerHost)
    // 运行程序
    core.Run(x.SerHost)
}
