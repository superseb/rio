package create

import (
	"fmt"

	"github.com/docker/cli/opts"
	"github.com/rancher/rio/cli/pkg/waiter"
	"github.com/rancher/rio/cli/server"
	"github.com/rancher/rio/types/client/rio/v1beta1"
	"github.com/urfave/cli"
)

type Create struct {
	AddHost              []string          `desc:"Add a custom host-to-IP mapping (host:ip)"`
	CapAdd               []string          `desc:"Add Linux capabilities"`
	CapDrop              []string          `desc:"Drop Linux capabilities"`
	Cidfile              string            `desc:"Write the container ID to the file"`
	Config               []string          `desc:"Configs to expose to the service (format: name:target)"`
	Cpus                 string            `desc:"Number of CPUs"`
	Device               []string          `desc:"Add a host device to the container"`
	Dns                  []string          `desc:"Set custom DNS servers"`
	DnsOption            []string          `desc:"Set DNS options"`
	DnsSearch            []string          `desc:"Set custom DNS search domains"`
	Entrypoint           []string          `desc:"Overwrite the default ENTRYPOINT of the image"`
	E_Env                []string          `desc:"Set environment variables"`
	EnvFile              []string          `desc:"Read in a file of environment variables"`
	GroupAdd             []string          `desc:"Add additional groups to join"`
	HealthCmd            string            `desc:"Command to run to check health"`
	HealthInterval       string            `desc:"Time between running the check (ms|s|m|h)" default:"0s"`
	HealthRetries        int               `desc:"Consecutive failures needed to report unhealthy"`
	HealthRecoverRetries int               `desc:"Consecutive failures needed to report healthy"`
	HealthStartPeriod    string            `desc:"Start period for the container to initialize before starting health-retries countdown (ms|s|m|h)" default:"0s"`
	HealthTimeout        string            `desc:"Maximum time to allow one check to run (ms|s|m|h)" default:"0s"`
	Hostname             string            `desc:"Container host name"`
	Init                 bool              `desc:"Run an init inside the container that forwards signals and reaps processes"`
	I_Interactive        bool              `desc:"Keep STDIN open even if not attached"`
	Ipc                  string            `desc:"IPC mode to use"`
	L_Label              map[string]string `desc:"Set meta data on a container"`
	LabelFile            []string          `desc:"Read in a line delimited file of labels"`
	M_Memory             string            `desc:"Memory limit (format: <number>[<unit>], where unit = b, k, m or g)"`
	MemoryReservation    string            `desc:"Memory soft limit (format: <number>[<unit>], where unit = b, k, m or g)"`
	N_Name               string            `desc:"Assign a name to the container"`
	Net_Network          string            `desc:"Connect a container to a network" default:"default"`
	NoHealthcheck        bool              `desc:"Disable any container-specified HEALTHCHECK"`
	Pid                  string            `desc:"PID namespace to use"`
	Privileged           bool              `desc:"Give extended privileges to this container"`
	P_Publish            []string          `desc:"Publish a container's port(s) externally"`
	ReadOnly             bool              `desc:"Mount the container's root filesystem as read only"`
	Restart              string            `desc:"Restart policy to apply when a container exits" default:"always"`
	SecurityOpt          []string          `desc:"Security Options"`
	StopTimeout          string            `desc:"Timeout (in seconds) to stop a container"`
	Tmpfs                []string          `desc:"Mount a tmpfs directory"`
	T_Tty                bool              `desc:"Allocate a pseudo-TTY"`
	U_User               string            `desc:"Username or UID (format: <name|uid>[:<group|gid>])"`
	V_Volume             []string          `desc:"Bind mount a volume"`
	VolumeDriver         string            `desc:"Optional volume driver for the container"`
	VolumesFrom          []string          `desc:"Mount volumes from the specified container(s)"`
	W_Workdir            string            `desc:"Working directory inside the container"`
}

func (c *Create) Run(app *cli.Context) error {
	return c.RunCallback(app, func(s *client.Service) *client.Service {
		return s
	})
}

func (c *Create) RunCallback(app *cli.Context, cb func(service *client.Service) *client.Service) error {
	var err error

	service, err := c.ToService(app.Args())
	if err != nil {
		return err
	}

	ctx, err := server.NewContext(app)
	if err != nil {
		return err
	}

	service.SpaceID, service.StackID, service.Name, err = ctx.ResolveSpaceStackName(service.Name)
	if err != nil {
		return err
	}

	service = cb(service)

	s, err := ctx.Client.Service.Create(service)
	if err != nil {
		return err
	}

	return waiter.WaitFor(app, s.ID)
}

func (c *Create) ToService(args []string) (*client.Service, error) {
	var (
		err error
	)

	if len(args) == 0 {
		return nil, fmt.Errorf("at least one (1) argument is required")
	}

	service := &client.Service{
		ExtraHosts:          c.AddHost,
		CapAdd:              c.CapAdd,
		CapDrop:             c.CapDrop,
		Command:             args[1:],
		CPUs:                c.Cpus,
		DefaultVolumeDriver: c.VolumeDriver,
		DNS:                 c.Dns,
		DNSOptions:          c.DnsOption,
		DNSSearch:           c.DnsSearch,
		Entrypoint:          c.Entrypoint,
		Hostname:            c.Hostname,
		Init:                c.Init,
		Image:               args[0],
		OpenStdin:           c.I_Interactive,
		IpcMode:             c.Ipc,
		Labels:              c.L_Label,
		Name:                c.N_Name,
		NetworkMode:         c.Net_Network,
		PidMode:             c.Pid,
		Privileged:          c.Privileged,
		ReadonlyRootfs:      c.ReadOnly,
		RestartPolicy:       c.Restart,
		//StopSignal:          c.StopSignal,
		Tty:         c.T_Tty,
		User:        c.U_User,
		VolumesFrom: c.VolumesFrom,
		WorkingDir:  c.W_Workdir,
	}

	service.Volumes, err = ParseMounts(c.V_Volume)
	if err != nil {
		return nil, err
	}

	service.Devices, err = ParseDevices(c.Device)
	if err != nil {
		return nil, err
	}

	service.Configs, err = ParseConfigs(c.Config)
	if err != nil {
		return nil, err
	}

	service.Environment, err = opts.ReadKVEnvStrings(c.EnvFile, c.E_Env)
	if err != nil {
		return nil, err
	}

	service.Labels, err = parseLabels(c.LabelFile, service.Labels)
	if err != nil {
		return nil, err
	}

	if err := populateHealthCheck(c, service); err != nil {
		return nil, err
	}

	if err := populateMemory(c, service); err != nil {
		return nil, err
	}

	service.Tmpfs, err = ParseTmpfs(c.Tmpfs)
	if err != nil {
		return nil, err
	}

	service.PortBindings, err = ParsePorts(c.P_Publish)
	if err != nil {
		return nil, err
	}

	return service, nil
}
