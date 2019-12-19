package core

// Command é a struct responsável por guardar informações referentes a um comando
type Command struct {
	Cmd         string `json:"command"`
	Description string `json:"description"`
	Usage       string `json:"usage"`
	Lint        string `json:"lint"`
	IsActive    bool   `json:"isActive"`
}

// Commands é a variável que guarda todos os comandos do BOT
var Commands []Command

// CreateCommands é a função que "cria" todos os comandos, essa função é chamada no main.go
func CreateCommands() {
	Commands = append(Commands, Command{
		Cmd:         canaryUpdate,
		Description: "Command that changes weights in Canary Deployment",
		Usage:       "@jeremias command `lb-id` `new-version-weight` `old-version-weight` `channel-to-send-alert (optional)`",
		Lint:        "`lb-id` LoadBalancer ID to be edited | `new-version-weight` Weight to new version on canary | `old-version-weight` Weight to old version on canary | `channel-to-send-alert` Channel code to send non-technical alert. Ex.: GHHG3S9L4",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         canaryUpTen,
		Description: "Command to add 10% to canary release of a load balancer",
		Usage:       "@jeremias command `lb-id` `channel-to-send-alert (optional)`",
		Lint:        "",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         canaryActivate,
		Description: "Command that actives the Canary Deployment in a specified Load Balancer",
		Usage:       "@jeremias command `*lb-id*`",
		Lint:        "The command removes all '#' of haproxy.cfg file | Will appear a select to you select a Load Balancer to enable canary",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         canaryDisable,
		Description: "Command that disable the Canary Deployment in a specified Load Balancer",
		Usage:       "@jeremias command `*lb-id*`",
		Lint:        "The command add '#' on start of all lines of the haproxy.cfg file | Will appear a select to you select a Load Balancer to enable canary",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         canaryInfo,
		Description: "Command that returns a haproxy.cfg of a specified Load Balancer",
		Usage:       "@jeremias command",
		Lint:        "The command get haproxy.cfg body and send for message | Will appear a select to you select a Load Balancer to get info",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         haproxyList,
		Description: "Command that brings ID list Environment Load Balancers Name",
		Usage:       "@jeremias command",
		Lint:        "",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         logsContainer,
		Description: "Command responsible for returning the logs of the specified container until the action is triggered",
		Usage:       "@jeremias command",
		Lint:        "Will appear a select, where be selected the container to get logs",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         restartContainer,
		Description: "Command responsible for restarting specified container",
		Usage:       "@jeremias command",
		Lint:        "Will appear a select, where be selected the container to be restarted",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         getServiceInfo,
		Description: "Command that brings information about a service that will be specified",
		Usage:       "@jeremias command",
		Lint:        "Will appear a select, where be selected the container",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         upgradeService,
		Description: "Command that will make an upgrade of a service, changing its image according to which it is passed as parameter",
		Usage:       "@jeremias command `service-id` `new-image`",
		Lint:        "In `service-id` put the id of the service which you need to send a new image and in `new-image` put the name of image to be sended",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         listService,
		Description: "Command that brings an ID list | Environment Services Name",
		Usage:       "@jeremias command",
		Lint:        "The returned format is something like ID: service-id | Name: service-name",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         startService,
		Description: "Command to activate services",
		Usage:       "@jeremias command `service-id`",
		Lint:        "",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         stopService,
		Description: "Command to deactivate a services",
		Usage:       "@jeremias command `service-id`",
		Lint:        "",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         checkServiceHealth,
		Description: "Command used to check health of one service",
		Usage:       "@jeremias command `stackName/serviceName` `channel-to-send-alert`",
		Lint:        "Put the Rancher Stack Name and Service Name on parameters, don't forget the '/'",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         taskAddByKeyword,
		Description: "Command used to add tasks using a keyword",
		Usage:       "@jeremias command `keyword1,keyword2` `channelToSendAlert` `deleteContainer?`",
		Lint:        "",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         statusService,
		Description: "Command used to check status of service",
		Usage:       "@jeremias command `stackName/serviceName` `channel-to-send-alert`",
		Lint:        "Put the Rancher Stack Name and Service Name on parameters, don't forget the '/'",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         removeServiceCheck,
		Description: "Command remove automatic service check",
		Usage:       "@jeremias command `task-ID`",
		Lint:        "Put the stackName and serviceName that you have informed on call to check. Task ID can be recovered with `list-task` command",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         listAllRunningTasks,
		Description: "Command list all running tasks at the moment of called",
		Usage:       "@jeremias command",
		Lint:        "Return a list with ID of running tasks",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         listAllEnvironments,
		Description: "Command list all environments of selected Rancher",
		Usage:       "@jeremias command",
		Lint:        "Return a list of all environments",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         selectEnvironment,
		Description: "Command to set a environment to next requests on Rancher",
		Usage:       "@jeremias command `environment-name`",
		Lint:        "The environment name can be recovered with environment-list command",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         envCleanupMachines,
		Description: "Command to cleanup machines of one environment on Rancher",
		Usage:       "@jeremias command",
		Lint:        "This command cleans disconnected machines from environment",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         selectRancher,
		Description: "Command sets the selected Rancher, to next requests",
		Usage:       "@jeremias command `rancher-name`",
		Lint:        "Receives Rancher name that has registered on database",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         listRancher,
		Description: "Command to list all registered Ranchers on database",
		Usage:       "@jeremias command",
		Lint:        "Returns `name, url and access key of all ranchers` (not returns secret key for security)",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         containerList,
		Description: "Command to list containers",
		Usage:       "@jeremias command",
		Lint:        "Returns `host, ID and name of all containers`",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         commands,
		Description: "Command responsible for displaying the commands that are available in BOT",
		Usage:       "@jeremias command",
		Lint:        "",
		IsActive:    true,
	})
}
