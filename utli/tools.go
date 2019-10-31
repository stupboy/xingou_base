package utli

import (
    "crypto/hmac"
    "crypto/md5"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "errors"
    "fmt"
    "math/rand"
    "regexp"
    "sort"
    "strconv"
    "time"
)

func GetTypeString(val interface{}) string {
    return fmt.Sprintf("%T", val)
}

func Reverse(x []interface{}) []interface{} {
    for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
        x[i], x[j] = x[j], x[i]
    }
    return x
}

//微信加密签名排序
func MapToString(data map[string]interface{}) (string, error) {
    var newMap = make([]string, 0)
    var str string
    var err error
    var temp string
    for k := range data {
        newMap = append(newMap, k)
    }
    sort.Strings(newMap)
    for _, v := range newMap {
        d := data[v]
        if len(str) == 0 {
            temp = ""
        } else {
            temp = "&"
        }
        switch d.(type) {
        case string:
            str += temp + v + "=" + d.(string)
        case int:
            str += temp + v + "=" + strconv.Itoa(d.(int))
        default:
            err = errors.New("未知类型错误")
            break
        }
    }
    return str, err
}

func Md5Encrypte(str string) string {
    h := md5.New()
    h.Write([]byte(str))
    return hex.EncodeToString(h.Sum(nil))
}

func ComputeHmacSha256(message string, secret string) string {
    key := []byte(secret)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(message))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
func Float64ToInt(arr map[string]interface{}, apiDoc interface{}) map[string]interface{} {
    newArr := make(map[string]interface{})
    for key, value := range arr {
        if apiDoc.(map[string]interface{})["param"].(map[string]interface{})[key] != nil {
            if apiDoc.(map[string]interface{})["param"].(map[string]interface{})[key].(map[string]interface{})["type"].(string) != "float" {
                switch value.(type) {
                case float64:
                    newArr[key] = strconv.Itoa(int(value.(float64)))
                default:
                    newArr[key] = value
                }
            } else {
                switch value.(type) {
                case string:
                    newArr[key], _ = strconv.ParseFloat(value.(string), 64)
                default:
                    newArr[key] = value
                }
            }
        }
    }
    return newArr
}

func CheckUrl(url string) bool {
    reg := regexp.MustCompile(`admin/[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
    if reg.FindAllString(url, -1) != nil {
        return true
    }
    return false
}

func GetRandomString(length, strType int) string {
    //strType 1大写英文+数字 2小写英文+数字
    var str string
    switch strType {
    case 1:
        str = "0123456789"
    case 2:
        str = "abcdefghijklmnopqrstuvwxyz"
    case 3:
        str = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    case 4:
        str = "0123456789abcdefghijklmnopqrstuvwxyz"
    case 5:
        str = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    case 6:
        str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    case 7:
        str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    }
    bytes := []byte(str)
    var result []byte
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < length; i++ {
        result = append(result, bytes[r.Intn(len(bytes))])
    }
    return string(result)
}

func GetCode(val int, flag int) string {
    var code string
    var str int
    if flag != 3 {
        sourceString := "598713624"
        num := val * 2
        for num > 0 {
            mod := num % 9
            num = (num - mod) / 9
            code = sourceString[mod:mod+1] + code
        }
        str, _ = strconv.Atoi(code)
    } else {
        str = val
    }
    return strconv.Itoa(str + 100000)
}

func Unique(m []string) []string {
    s := make([]string, 0)
    sMap := make(map[string]string)
    for _, value := range m {
        if value != "" {
            //计算map长度
            length := len(sMap)
            sMap[value] = "1"
            //比较map长度, 如果map长度不相等， 说明key不存在
            if len(sMap) != length {
                s = append(s, value)
            }
        }
    }
    return s
}

func InArray(arr []string, str string) bool {
    ret := false
    for _, val := range arr {
        if str == val {
            ret = true
        }
    }
    return ret
}

func InArrayInt(arr []int, str int) bool {
    ret := false
    for _, val := range arr {
        if str == val {
            ret = true
        }
    }
    return ret
}

func RandNum() string {
    return strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(899999) + 100000)
}
