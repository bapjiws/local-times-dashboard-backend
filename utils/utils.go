package utils

func PanicOnError(e error) {
	if e != nil {
		panic(e.Error())
	}
}
