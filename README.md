## rtbot.go

A Twitter bot that tweets TechCrunch Japan articles that have garnered more than 100RTs and beyond. This bot checks the retweet counts for articles within three days every 10 minutes and it generates a new tweet every time an article gets additional 50RTs like so: 100, 150, 200, 250, 300...

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
- add web ui to view the current stats
- add weekly / monthly top 5 tweets
- take config file path as a command line option
- create twitter buffer goroutine to optimize the timing
- logging
- follow everyone who retweets or favs our tweets, and unfollow those who don't follow us back in a week
- retweet all the tweets with an opinion added to the original tweet
- generalize the article extraction process to use this bot for other sites

## License

Copyright (c) 2015 Ken Nishimura
This software is released under the MIT License
http://rem.mit-license.org

