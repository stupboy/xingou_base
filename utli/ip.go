package utli

import (
    "bytes"
    "strconv"
)

func IpToInt(ip []byte) int {
    var ipInt int = 0
    var pos uint = 24
    for _, ipSeg := range ip {
        tempInt := int(ipSeg)
        tempInt = tempInt << pos
        ipInt = ipInt | tempInt
        pos -= 8
    }
    return ipInt
}

func IntToIp(ipInt int) string {
    ipSegs := make([]string, 4)
    var len1 int = len(ipSegs)
    buffer := bytes.NewBufferString("")
    for i := 0; i < len1; i++ {
        tempInt := ipInt & 0xFF
        ipSegs[len1-i-1] = strconv.Itoa(tempInt)
        ipInt = ipInt >> 8
    }
    for i := 0; i < len1; i++ {
        buffer.WriteString(ipSegs[i])
        if i < len1-1 {
            buffer.WriteString(".")
        }
    }
    return buffer.String()
}
