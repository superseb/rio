package lookup

import (
	"github.com/rancher/rio/cli/pkg/kv"
)

type ParsedService struct {
	StackName   string
	ServiceName string
	Revision    string
}

func ParseServiceName(serviceName string) ParsedService {
	var result ParsedService
	serviceName, result.Revision = kv.Split(serviceName, ":")
	result.StackName, result.ServiceName = kv.Split(serviceName, "/")
	if result.ServiceName == "" {
		result.ServiceName = result.StackName
		result.StackName = "default"
	}
	return result
}

func (p ParsedService) Latest() ParsedService {
	return p.SetRevision("latest")
}

func (p ParsedService) SetRevision(rev string) ParsedService {
	p.Revision = rev
	return p
}

func (p ParsedService) String() string {
	result := ""
	if p.StackName != "" && p.StackName != "default" {
		result = p.StackName + "/"
	}
	result += p.ServiceName

	if p.Revision != "" && p.Revision != "latest" {
		result += ":" + p.Revision
	}

	return result
}
