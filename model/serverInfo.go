package model

//MsInfo ...
type MsInfo struct {
	Host string
	UUID string
}

//APIInfo ...
type APIInfo struct {
	APIName      string
	APIMethods   string
	TokenMethods string
}

//RequestModel ...
type RequestModel struct {
	MsName       string
	APIName      string
	APIMethods   string
	TokenMethods string
	MSHost       string
}

//ComResult ..
type ComResult struct {
	StatusCode int
	Code       int
	Msg        string
	Info       interface{}
}
