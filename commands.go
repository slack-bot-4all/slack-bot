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
		Cmd:         "haproxy-update",
		Description: "Comando que faz alteração nos pesos do Canary Deployment",
		Usage:       "@bot comando `id-lb` `porc-new` `porc-old`",
		Lint:        "`id-lb` ID do Load Balancer a ser editado | `porc-new` Porcentagem que será adicionada na nova versão | `porc-old` Porcentagem que será adicionada na antiga versão",
	},
	{
		Cmd:         "lb-list",
		Description: "Comando que trás a lista de ID | Nome dos Load Balancers do Environment",
		Usage:       "@bot comando",
		Lint:        "",
	},
	{
		Cmd:         "logs-container",
		Description: "Comando que trará um arquivo com os arquivos de logs do container selecionado",
		Usage:       "@bot comando",
		Lint:        "Aparecerá uma caixa de seleção, onde será selecionado o container que preferir",
	},
	{
		Cmd:         "restart-container",
		Description: "Comando que reinicia o container selecionado",
		Usage:       "@bot comando",
		Lint:        "Aparecerá uma caixa de seleção, onde será selecionado o container a ser restartado",
	},
	{
		Cmd:         "info-service",
		Description: "Comando que busca informações sobre o serviço selecionado",
		Usage:       "@bot comando",
		Lint:        "Aparecerá uma caixa de seleção, onde será selecionado o serviço a ser buscado",
	},
	{
		Cmd:         "comandos",
		Description: "Comando responsável por mostrar os comandos que estão disponíveis no BOT",
		Usage:       "@bot comando",
		Lint:        "",
	},
}
