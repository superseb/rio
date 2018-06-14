package server

type Server struct {
	P_HttpsListenPort  int    `desc:"HTTPS listen port" default:"7443"`
	L_HttpListenPort   int    `desc:"HTTP listen port" default:"7080"`
	D_DataDir          string `desc:"Folder to hold state" default:"${HOME}/.rancher/rio/server"`
	DisableControllers bool   `desc:"Don't run controllers (only useful for rio development)"`
}
