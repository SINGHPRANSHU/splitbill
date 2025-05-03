package handler

type LoggerInterface interface {
	Debug(v ...interface{})
	Fatal(v ...interface{})
	Println(v ...interface{})
}
