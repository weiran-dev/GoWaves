package util

func HandleError(err error) {
	if err != nil {
		panic("出现致命错误：" + err.Error())
	}
}