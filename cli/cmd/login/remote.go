package login

import (
	"github.com/rancher/rio/cli/pkg/up/questions"
	"github.com/rancher/rio/pkg/clientaccess"
)

func (l *Login) remote(tempFile string) error {
	var err error

	if l.S_Server == "" {
		l.S_Server, err = questions.Prompt("Rio server URL: ", "")
		if err != nil {
			return err
		}
	}

	if l.T_Token == "" {
		l.T_Token, err = questions.Prompt("Authentication token: ", "")
		if err != nil {
			return err
		}
	}

	return clientaccess.AccessInfoToKubeConfig(tempFile, l.S_Server, l.T_Token)
}
