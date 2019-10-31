package utli

import (
    "io/ioutil"
    "net/http"
    "strings"
)

func HttpGet(url string, param map[string]string) ([]byte, error) {
    var paramString string
    for key, val := range param {
        if paramString == "" {
            paramString = "?" + key + "=" + val
        } else {
            paramString = paramString + "&" + key + "=" + val
        }
    }
    url = url + paramString
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)
}

func HttpPostJson(url string, param map[string]string, data []byte) ([]byte, error) {
    var paramString string
    for key, val := range param {
        if paramString == "" {
            paramString = "?" + key + "=" + val
        } else {
            paramString = paramString + "&" + key + "=" + val
        }
    }
    url = url + paramString
    resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)
}

func HttpPost(url string, param map[string]string) ([]byte, error) {
    var paramString string
    for key, val := range param {
        if paramString == "" {
            paramString = key + "=" + val
        } else {
            paramString = paramString + "&" + key + "=" + val
        }
    }
    resp, err := http.Post(url,
        "x-www-form-urlencoded",
        strings.NewReader(paramString))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)
}

func HttpDo(url string, method string, param map[string]string, data map[string]string) ([]byte, error) {
    client := &http.Client{}
    var paramString string
    for key, val := range param {
        if paramString == "" {
            paramString = key + "=" + val
        } else {
            paramString = paramString + "&" + key + "=" + val
        }
    }
    req, err := http.NewRequest(method, url, strings.NewReader(paramString))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    for key, val := range data {
        req.Header.Set(key, val)
    }
    // req.Header.Set("Cookie", "name=anny")
    resp, err := client.Do(req)
    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)
}
