package main

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
		Description: "Comando que faz alteração nos pesos do Canary Deployment",
		Usage:       "@bot comando `id-lb` `porc-new` `porc-old`",
		Lint:        "`id-lb` ID do Load Balancer a ser editado | `porc-new` Porcentagem que será adicionada na nova versão | `porc-old` Porcentagem que será adicionada na antiga versão",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         canaryActivate,
		Description: "Comando que ativa o Canary Deployment",
		Usage:       "@bot comando `*id-lb*`",
		Lint:        "O comando tira todos os '#' que tem no arquivo haproxy.cfg | Aparecerá um select onde você selecionará o Load Balancer ou você pode enviar o ID do LB por parâmetro",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         canaryDisable,
		Description: "Comando que desativa o Canary Deployment",
		Usage:       "@bot comando `*id-lb*`",
		Lint:        "O comando adiciona um '#' no início de todas as linhas que tem no arquivo haproxy.cfg | Aparecerá um select onde você selecionará o Load Balancer ou você pode enviar o ID do LB por parâmetro",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         canaryInfo,
		Description: "Comando que trás o haproxy.cfg do Load Balancer informado, com propósito de trazer as informações do Canary Deployment",
		Usage:       "@bot comando",
		Lint:        "O comando busca o haproxy.cfg e apenas envia como mensagem | Aparecerá um select onde você selecionará o Load Balancer",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         haproxyList,
		Description: "Comando que trás a lista de ID | Nome dos Load Balancers do Environment",
		Usage:       "@bot comando",
		Lint:        "",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         logsContainer,
		Description: "Comando que trará um arquivo com os arquivos de logs do container selecionado",
		Usage:       "@bot comando",
		Lint:        "Aparecerá uma caixa de seleção, onde será selecionado o container que preferir",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         restartContainer,
		Description: "Comando que reinicia o container selecionado",
		Usage:       "@bot comando",
		Lint:        "Aparecerá uma caixa de seleção, onde será selecionado o container a ser restartado",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         getServiceInfo,
		Description: "Comando que busca informações sobre o serviço selecionado",
		Usage:       "@bot comando",
		Lint:        "Aparecerá uma caixa de seleção, onde será selecionado o serviço a ser buscado",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         upgradeService,
		Description: "O comando faz o upgrade de um serviço mudando apenas sua imagem",
		Usage:       "@bot comando `id-serviço` `nova-imagem`",
		Lint:        "Em `id-serviço` coloque o ID referente ao serviço que você quer enviar a nova imagem e em `nova-imagem` coloque o nome da imagem a ser enviada",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         listService,
		Description: "O comando lista todos os serviços disponíveis no Environment, listando de forma resumida apenas ID e Nome",
		Usage:       "@bot comando",
		Lint:        "O formato de retorno será algo como ID: id-serviço | Nome: nome-serviço",
		IsActive:    true,
	})

	Commands = append(Commands, Command{
		Cmd:         comandos,
		Description: "Comando responsável por mostrar os comandos que estão disponíveis no BOT",
		Usage:       "@bot comando",
		Lint:        "",
		IsActive:    true,
	})
}
