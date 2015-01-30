## rtbot.go

A Twitter bot that tweets TechCrunch Japan articles that have garnered more than 100RTs and beyond. Every ten minutes, this bot checks the retweet counts for articles within three days and generates a new tweet every time an article gets additional 10RTs like so: 100, 200, 300, 400 and so on.

~~~
$ ./rtbot --dry-run=false -d 3 -interval=10m

--dry-run | -d | -interval

--dry-run: suppress tweeting, just show tweets on stdout
-d       : days to look back and check
-interval: polling interval
~~~

conf.json file should look like:
~~~
{
    "consumerkey": "xxxyyyzzz111222333",
    "consumersecret": "xxxyyyzzz111222333",
    "accesstoken": "xxxyyyzzz111222333",
    "tokensecret": "xxxyyyzzz111222333"
}
~~~

## Dependency

Golang twitter library:
https://github.com/ChimeraCoder/anaconda

## Todo:
- add weekly / monthly top 5 tweets
- take config file path as a command line option
- create tweets buffer goroutine to optimize the timing of tweet
- follow everyone who retweets or favs our tweets, and unfollow those who don't follow us back in a week
- retweet all the tweets with an opinion added to the original tweet
- generalize the article extraction process to use this bot for other sites

## License

Copyright (c) 2015 Ken Nishimura
This software is released under the MIT License
http://rem.mit-license.org
