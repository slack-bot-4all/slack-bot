(not updated readme)

# Slack-bot for Rancher

- [How to Works?](#How-to-Works)
- [How to Use?](#How-to-use)
- [Available Commands](#Available-Commands)
- [Scheduling Commands](#scheduling-commands)
- [Contribution](#Contribution)
- [Adding New Commands](#Adding-New-Commands)

The ***SLfR*** (Slack-bot for Rancher), is an application responsible for task automation in Rancher 1.6, using the Rancher and Slack API.

The bot's intent is to facilitate tasks that are common to Rancher users, such as:

    - Check container health
    - Can restart containers
    - Getlogs
    - Upgrades
    - Canary Releases

## How to Works?

To better understand how a bot works inside the slack, learn more here [Slack-bot](https://api.slack.com/bot-users) .

Two of the features that the BOT comes from are the ***@rancher_bot restart-container*** and ***@rancher_bot logs-container***. When calling the bot by passing any of the features, the return will look something like this:

![restart_container](images/restart-container.png)

The bot will return you a list of containers running in your enviromment, which will be defined in the file **.env**. 
The application is dockerized in order to speed up the execution and facilitate its scalability, if need be.

Given the bot's operation in the slack, we raised an API to make your URI available to the slack, which is where the interpretations of the ***Actions*** registered to the bot will happen.

## How to use?

First, let's go to the <a href="https://api.slack.com/apps" target="_blank"> Slack APPS </a> site, you'll fall into a page something like this:

![your_apps](images/your-apps.PNG)

Let's click **Create New App**, it will open a window asking **App Name** that you will put the name of your application, and **Development Slack Workspace** that you will place where your app will work, which in the case, it will be on your Slack server. After that, click **Create App**

After that, you will be redirected to the menu of your application, go to the option **OAuth & Permissions**, in Scopes> Select Permission Scopes you must select the permissions that your application will have, we recommend that you put the **Administer the workspace**, so that the application has access and permission to everything within Slack, but if your server is more careful, put the permissions as you prefer. Remember that the permissions you put in may influence the operation of the BOT.

After putting the permission, the part of **Scopes** will look something like this:

![permissions](images/permissions.PNG)

Just click **Save Changes** to save what was done.

After that, go to the option **Bot Users** and click **Add a Bot User** to add a new BOT to your application, customize the name that will be displayed and the user, and click **Save Changes**

After that, go back to **Basic Information** in the menu and look for **Install your app to your workspace**, click on this option, see that it will show some information next to the **Install App to Workspace** button, click button. When you click the button to install the application on your server, it will open a page asking you to authorize the application, click **Authorize**

To use ***SLfR*** is simple, you will first need to download the source code (by cloning this repository) console slack-bot @ pc: ~ $ git clone https://github.com/slack-bot-4all/slack-bot.git
(I.e.
After that, change the ```.env``` file by adding the Rancher 1.6, BOT and HTTP port information that will run the API
```properties
RANCHER_ACCESS_KEY=<RANCHER_API_ACCESS_KEY>
RANCHER_SECRET_KEY=<RANCHER_API_SECRET_KEY>
RANCHER_BASE_URL=<API_BASE_URL> Ex.: http://yourdomain.ip:8080/v1/projects
RANCHER_PROJECT_ID=<ENVIRONMENT_ID>
SLACK_BOT_TOKEN=<API_SLACK_ACCESS_TOKEN>
SLACK_BOT_ID=<BOT_ID>
SLACK_BOT_CHANNEL=<CHANNEL_WHERE_THE_BOT_LISTEN_COMMANDS>
SLACK_BOT_VERIFICATION_TOKEN=<BOT_VERIFICATION_TOKEN>
HTTP_PORT=<HTTP_PORT>
```

**Note: To get the BOT ID, you will need to first leave it blank and run the application (which will be taught below), you will get the BOT ID in the application logs, as in the image below.**

![id-bot](images/id-bot.PNG)

With the ```.env``` file changed, you will need to decide how to run, whether to run on the machine where you are connected, or if you want to ***dockerizer***. If you want to run directly on the machine, just run the ```.go``` files, as follows:
```console
slack-bot@pc:~$ go run *.go
```
If you want ***dockerizar***, just do the Docker image build, our ```Dockerfile``` is ready to be build:
```console
slack-bot@pc:~$ docker build -t usuario/nome-da-imagem:versao .
```
And after that, just give the ***docker run*** in your already-built image:
```console
slack-bot@pc:~$ docker run -d -p PORT_HTTP:PORT_HTTP -e "FILE=.env" user/image-name:version
```
Remember to externalize the HTTP port you set in ```.env```, so that the Slack API can access the URL.

**Done, now the BOT is already running, just check the Slack channel you have set for him to listen to the messages if he sent the message telling you it's online, and just use it! :blush:**

## Available Commands

| Command | Description |
| ------- | --------- |
| `restart-container` | *Command responsible for restarting specified container* |
| `logs-container` | *Command responsible for returning the logs of the specified container until the action is triggered* |
| `update-canary` | *Command that changes weights in Canary Deployment* |
| `enable-canary` | *Command that actives the Canary Deployment in a specified Load Balancer* |
| `disable-canary` | *Command that disable the Canary Deployment in a specified Load Balancer* |
| `info-canary` | *Command that returns a haproxy.cfg of a specified Load Balancer* |
| `list-lb` | *Command that brings ID list Environment Load Balancers Name* |
| `info-service` | *Command that brings information about a service that will be specified* |
| `upgrade-service` | *Command that will make an upgrade of a service, changing its image according to which it is passed as parameter* |
| `list-service` | *Command that brings an ID list \| Environment Services Name* |
| `start-service` | *Command that start one service* |
| `stop-service` | *Command that stop one service* |
| `check-service` | *Command that check the health of one service* |
| `stop-check` | *Command that stop the one already created check* |
| `commands` | *Command responsible for displaying the commands that are available in BOT* |

## Scheduling Commands
Our BOT is adapted to receive a "reminder" messages, that way, the BOT processes the message and take the command. A simple usage of [Slack Reminder](https://get.slack.help/hc/en-us/articles/208423427-Definir-um-lembrete) is:
```
/remind channel-to-send-a-message "My Command" time
```
A simple example of usage is:
```
/remind #general "@jeremias update-canary 1s30 90 10" at 11:00pm
```

## Contribution
We are fully open to contri- butions. This is an **Open Source** project, so whatever you have to add in our project, just add and do the pull request.

## Adding New Commands
If it is necessary to add new commands, simply add the constant in `slack.go`, in the group of global constants
```golang
const (
    yourCommand = "SlacklikeCommand"
)
```
After that, still in `slack.go`, look for the `handleMessageEvent` function to go to the end of the function, see that it will have a condition structure chair (IfElse), add another `else if () {}` with the following rules:
```golang
else if strings.HasPrefix(message, yourCommand) {
    funcProcessCommand()
}
```

After that, add your command to `commands.go`, inside the **slice** called `Commands`. *Note: This step is optional, if not put, your command will not appear in the command list, however, it will work*.

```golang
Commands = append(Commands, Command{
		Cmd:         yourCommand,
		Description: "Description of your command, explaining what it is for",
		Usage:       "As your command will be used (we recommend that you refer to the command as 'command', because when the command listing method is called, it will be replaced by the command itself)",
		Lint:        "If your command receives arguments or you want to leave any tips on the command, put it here",
		IsActive:    true , // Keep true. Its a possible feature.
            })
```
