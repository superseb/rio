package agent

type Agent struct {
	T_Token   string `desc:"Token to use for authentication" env:"RIO_TOKEN"`
	S_Server  string `desc:"Server to connect to" env:"RIO_URL"`
	D_DataDir string `desc:"Folder to hold state" default:"${HOME}/.rancher/rio/agent"`
}
