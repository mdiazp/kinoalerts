# Prerequisites
- Telegram Bot Token

# Installation on Linux
- Open the file env-examples.sh and set your correct bot token.
- Execute the command: ./run-on-linux.sh

# Configuration
## Environment config:
- TELEGRAM_BOT_TOKEN: The token of the bot you will use to notificate
- INTERVAL_TIME_TO_QUERY: How many seconds the system delays between calls to the resources
- INTERVAL_TIME_TO_NOTIFICATE: How many seconds the system notificate about some alert already founded

## Resources config:
The resources.json file has this structue:
```
{
    "resources": [
        {
            "name": "TuEnvio => Pinar => Alimentos Refrigerados",
            "chatID": "-1001283379000",
            "url": "https://www.tuenvio.cu/pinar/Products?depPid=46081",
            "alerts": [
                {
                    "name": "Pollo",
                    "regexPatterns": ["pollo"]
                },
                {
                    "name": "Carne de Res",
                    "regexPatterns": [" res "]
                },
                {
                    "name": "Salchichon",
                    "regexPatterns": ["salchichon"]
                }
            ]
        },
    ]
} 
```
You can configurate several resources.

### Resource:
A resource is a link where you want the system looks for something. In the url property you need to put this link. The property chatID is the telegram chat id to notificate (Example: A telegram group that your bot is a member).
Every resource has an array of alerts.

### Alert:
An alert is some info that you want to look if it exists in the resource.
The name is the text in the telegram notification to identify this alert, and
the regexPatterns is an array of the regular expressionsto to match.

# Practice Use Case
The online shop https://www.tuenvio.cu offers several products, and you wants to be notified when some products on the section [Pinar del Rio => Alimentos Refrigerados] are availables.

## Steps
1. Create the telegram bot and get the bot_token
2. Create the telegram group where you want the bot to notifiy
3. Add the bot as a member of the group
4. Get the telegram group chatID
5. Set the env-example.sh file with the bot_token correctly configured
6. Set the resource.json file with the correct configuration
```
{
    "resources": [
        {
            "name": "TuEnvio => Pinar => Alimentos Refrigerados",
            "chatID": "chatGroupID",
            "url": "https://www.tuenvio.cu/pinar/Products?depPid=46081",
            "alerts": [
                {
                    "name": "Pollo",
                    "regexPatterns": ["pollo"]
                }
            ]
        },
    ]
} 
```



