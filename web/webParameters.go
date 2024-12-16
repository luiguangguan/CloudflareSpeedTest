package web

type SubmitData struct {
	Password string
	Action   string
	Content  string
}

type ConfigData struct {
	Password string
	Content  string
}
type TestHttpConnectData struct {
	Password     string
	IpText       string
	TestDownload bool
}
type Pwd struct {
	Password string
}

type EditPassword struct {
	OldPwd  string
	NewPwd1 string
	NewPwd2 string
}
