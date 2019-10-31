package utli

import (
    "bufio"
    "encoding/json"
    "io"
    "os"
)

func GetApiDocJson() map[string]interface{} {
    var filename = "doc.json"
    data := make(map[string]interface{})
    doc, _ := os.Open(filename)
    rd1 := bufio.NewReader(doc)
    line, err := rd1.ReadString('\n') //以'\n'为结束符读入一行
    if err == nil || err == io.EOF {
        err = json.Unmarshal([]byte(line), &data)
        if err == nil {
            return data
        }
    } else {
        panic("配置文件不存在")
    }
    return data
}
