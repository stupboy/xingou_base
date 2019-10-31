package validate

import (
    "strconv"
    "strings"
)

func LenValidate(val string, argusStr ...string) (bool) {
    if len(argusStr) == 0 {
        return false
    }
    argus := strings.Split(argusStr[0], "-")
    if len(argus) > 2 || len(argus) < 1 {
        return false
    }
    max := 0
    min := 0
    min, _ = strconv.Atoi(argus[0])
    if len(argus) == 2 {
        max, _ = strconv.Atoi(argus[1])
    }
    vlen := len(val)
    if max == 0 {
        if vlen == min {
            return true
        }
    }
    if max < min {
        return false
    }
    if min <= vlen && vlen <= max {
        return true
    }
    return false
}

func MaxValidate(val string, argusStr ...string) bool {
    valInt, _ := strconv.Atoi(val)
    if len(argusStr) == 0 {
        return false
    }
    max, _ := strconv.Atoi(argusStr[0])
    if valInt <= max {
        return true
    }
    return false
}

func MinValidate(val string, argusStr ...string) bool {
    valInt, _ := strconv.Atoi(val)
    if len(argusStr) == 0 {
        return false
    }
    max, _ := strconv.Atoi(argusStr[0])
    if valInt >= max {
        return true
    }
    return false
}

func RangValidate(val string, argusStr ...string) bool {
    if len(argusStr) == 0 {
        return false
    }
    argus := strings.Split(argusStr[0], "-")
    if len(argus) != 2 {
        return false
    }
    valInt, _ := strconv.Atoi(val)
    min, _ := strconv.Atoi(argus[0])
    max, _ := strconv.Atoi(argus[1])
    if valInt >= min && valInt <= max {
        return true
    }
    return false
}

func MobileValidate(val string, argusStr ...string) bool {
    if len(val) != 11 {
        return false
    }
    if val[0:1] != "1" {
        return false
    }
    if val[0:2] == "12" {
        return false
    }
    return true
}

func InValidate(val string, argusStr ...string) bool {
    if len(argusStr) == 0 {
        return false
    }
    argusArr := strings.Split(argusStr[0], "-")
    ret := false
    for _, v := range argusArr {
        if v == val {
            ret = true
            break
        }
    }
    return ret
}
