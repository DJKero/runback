package error

import "log"

// Checks error, logs it, and recovers if possible
func ErrorLog(e error) {
	if e != nil {
		log.Println("New Error:")
		log.Println(e)
		recover()
	}
}

func ErrorLogCustom(s string, e error) {
	if e != nil {
		log.Println(s)
		log.Println(e)
		recover()
	}
}
