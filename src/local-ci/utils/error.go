package utils

func CheckError(err error) {
	if nil != err {
		panic(err)
	}
}
