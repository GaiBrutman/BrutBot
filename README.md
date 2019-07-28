# BrutBot
A Discord bot that uses the Reddit API to fetch random posts of dogs, cats, memes and more


## Setup
first initialize a new Discod bot.  
After setting up the bot and receiving a token, you will need to create a config.json file. looking like this:
```
{
  "Token": [TOKEN],
  "BotPrefix": [PREFIX]
}
```
## Commands
A command will start with a prefix specified in the config.json file, followed by the command and its arguments:  
```[Prefix][Command] [Arg1] [Arg2] ...```

for example:
```$reddit dogpictures```

*the commands are case insensitive

**ping:** sends back "Pong".  
**help / h:** sends the bot's description and commanding help.  
**time / t:** sends the current time.  
**gopher / go:** sends a Golang gopher image.  
**dog / cat:** sends a random Reddit post of dogs/cats (respectively).  
**pewdiepie / pewds:** sends a random Reddit post the 'r/PewdiepieSubmissions' subreddit.  
**reddit / rd [Arg]:** sends a random Reddit post from a given subreddit.  
