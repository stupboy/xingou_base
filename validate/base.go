package validate

import (
    "fmt"
    "strconv"
    "strings"
)

// 请求参数验证
func RequestValidate(requestMap map[string]interface{}, ruleMap map[string]interface{}, ValidateMap map[string]interface{}) map[string]interface{} {
    data := make(map[string]interface{})
    // 不需要参数情况
    if _, ok := ruleMap["param"]; !ok {
        return nil
    }
    defer func() {
        if r := recover(); r != nil {
            panic(r)
        }
    }()
    // 根据注释获取需要参数
    for key, val := range ruleMap["param"].(map[string]interface{}) {
        valRule := val.(map[string]interface{})
        if _, ok := requestMap[key]; !ok && valRule["must"].(string) == "1" {
            panic(key + "-参数不存在")
        }
        if valRule["type"].(string) != "string" && valRule["type"].(string) != "json" && requestMap[key] == "" {
            delete(requestMap, key)
        }
        if _, ok := requestMap[key]; ok {
            switch valRule["type"].(string) {
            case "int":
                val, _ := strconv.Atoi(requestMap[key].(string))
                if fmt.Sprintf("%v", val) != fmt.Sprintf("%v", requestMap[key].(string)) {
                    panic("请输入正确数字" + key)
                }
            }
        }
        if _, ok := requestMap[key]; !ok {
            if valRule["value"].(string) != "none" {
                data[key] = valRule["value"].(string)
                continue
            }
        }
        if _, ok := requestMap[key]; !ok {
            continue
        }
        // TODO 其他验证规则
        paramRule := valRule["rule"]
        if paramRule != "none" {
            paramRules := strings.Split(paramRule.(string), "|")
            for _, r := range paramRules {
                a1 := strings.Split(r, ":")
                var ret bool
                if _, ok := ValidateMap[a1[0]]; !ok {
                    panic("验证规则不存在:" + a1[0])
                }
                if len(a1) == 1 {
                    ret = ValidateMap[a1[0]].(func(string, ...string) bool)(requestMap[key].(string))
                }
                if len(a1) == 2 {
                    ret = ValidateMap[a1[0]].(func(string, ...string) bool)(requestMap[key].(string), a1[1])
                }
                if !ret {
                    errStr := valRule["info"].(string) + key + "的验证规则不通过,规则值:" + r
                    panic(errStr)
                }
            }
        }
        data[key] = requestMap[key]
    }
    return data
}

// 返回参数验证
func ReturnValidate(returnMap map[string]interface{}, ruleMap map[string]interface{}) {
    // TODO 返回值的校验方法
    // 遍历返回值 检查字段名 是否确实  返回类型是否正确
}
