package main

// Command é a struct responsável por guardar informações referentes a um comando
type Command struct {
	Cmd         string
	Description string
	Usage       string
	Lint        string
}

// Commands é a variável que guarda todos os comandos do BOT
var Commands = []Command{
	{
		Cmd:         canaryUpdate,
		Description: "Comando que faz alteração nos pesos do Canary Deployment",
		Usage:       "@bot comando `id-lb` `porc-new` `porc-old`",
		Lint:        "`id-lb` ID do Load Balancer a ser editado | `porc-new` Porcentagem que será adicionada na nova versão | `porc-old` Porcentagem que será adicionada na antiga versão",
	},
	{
		Cmd:         canaryActivate,
		Description: "Comando que ativa o Canary Deployment",
		Usage:       "@bot comando `*id-lb*`",
		Lint:        "O comando tira todos os '#' que tem no arquivo haproxy.cfg | Aparecerá um select onde você selecionará o Load Balancer ou você pode enviar o ID do LB por parâmetro",
	},
	{
		Cmd:         canaryDisable,
		Description: "Comando que desativa o Canary Deployment",
		Usage:       "@bot comando `*id-lb*`",
		Lint:        "O comando adiciona um '#' no início de todas as linhas que tem no arquivo haproxy.cfg | Aparecerá um select onde você selecionará o Load Balancer ou você pode enviar o ID do LB por parâmetro",
	},
	{
		Cmd:         canaryInfo,
		Description: "Comando que trás o haproxy.cfg do Load Balancer informado, com propósito de trazer as informações do Canary Deployment",
		Usage:       "@bot comando",
		Lint:        "O comando busca o haproxy.cfg e apenas envia como mensagem | Aparecerá um select onde você selecionará o Load Balancer",
	},
	{
		Cmd:         haproxyList,
		Description: "Comando que trás a lista de ID | Nome dos Load Balancers do Environment",
		Usage:       "@bot comando",
		Lint:        "",
	},
	{
		Cmd:         logsContainer,
		Description: "Comando que trará um arquivo com os arquivos de logs do container selecionado",
		Usage:       "@bot comando",
		Lint:        "Aparecerá uma caixa de seleção, onde será selecionado o container que preferir",
	},
	{
		Cmd:         restartContainer,
		Description: "Comando que reinicia o container selecionado",
		Usage:       "@bot comando",
		Lint:        "Aparecerá uma caixa de seleção, onde será selecionado o container a ser restartado",
	},
	{
		Cmd:         getServiceInfo,
		Description: "Comando que busca informações sobre o serviço selecionado",
		Usage:       "@bot comando",
		Lint:        "Aparecerá uma caixa de seleção, onde será selecionado o serviço a ser buscado",
	},
	{
		Cmd:         upgradeService,
		Description: "O comando faz o upgrade de um serviço mudando apenas sua imagem",
		Usage:       "@bot comando `id-serviço` `nova-imagem`",
		Lint:        "Em `id-serviço` coloque o ID referente ao serviço que você quer enviar a nova imagem e em `nova-imagem` coloque o nome da imagem a ser enviada",
	},
	{
		Cmd:         listService,
		Description: "O comando lista todos os serviços disponíveis no Environment, listando de forma resumida apenas ID e Nome",
		Usage:       "@bot comando",
		Lint:        "O formato de retorno será algo como ID: id-serviço | Nome: nome-serviço",
	},
	{
		Cmd:         comandos,
		Description: "Comando responsável por mostrar os comandos que estão disponíveis no BOT",
		Usage:       "@bot comando",
		Lint:        "",
	},
}
