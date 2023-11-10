package e

var MsgFlags = map[int]string{
	Success:     "ok",
	Unknown:     "未知错误",
	ValidateErr: "请求参数错误",
	NotFoundErr: "资源不存在",
}

// 从错误码中获取错误提示
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Unknown]
}
