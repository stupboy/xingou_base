package utli

const ReturnSuccess = 1
const ReturnError = 0

const DEFAULTERR = 10000
const USERNOTEXIST = 10001
const (
    AuthInfoNotExist = -14
)

var ErrCodeMap = map[int]string{
    1:     "ok",
    -14:   "鉴权信息不存在",
    0:     "failed",
    10000: "未定义错误",
}
