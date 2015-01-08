package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"
)

const (
	// This Twitter API is actually an unofficial endpoint. We
	// should use the Streaming API instead.
	retweetAPI string = "http://urls.api.twitter.com/1/urls/count.json?url="
	datafile string = "articles.csv"
)

type Article struct {
	published time.Time
	url       string
	title     string
	retweet   int
}

var Config struct {
	// Defines the days to look back in the past
	days int
	// Polling duration
	interval time.Duration
	// Suppress tweets in dry-run mode
	dryrun bool
}

func init() {
	flag.IntVar(&Config.days, "d", 3, "Days to look back")
	flag.DurationVar(&Config.interval, "interval", 10*time.Minute, "Polling interval")
	flag.BoolVar(&Config.dryrun, "dry-run", false, "Dry run mode")
	flag.Parse()

	if Config.dryrun {
		fmt.Println("Runnunig in dry-run mode.")
	}
}

func main() {
	client := NewTwitterClient("conf.json")

	// Check and store the current stats for the first time
	articles := withinDays(fetchArticles(3), Config.days)
	savecsv(articles)
	fmt.Println("Saved the CSV file.")

	go startServer()
	fmt.Println("Started the http server.")

	tick := time.Tick(Config.interval)
	for range tick {
		articles := withinDays(fetchArticles(3), Config.days)
		prevArticles := withinDays(loadcsv(), Config.days)
		tweetedUrls := createMap(prevArticles)

		for _, a := range articles {
			if TweetWorthy(a.retweet, tweetedUrls[a.url]) {
				msg := fmt.Sprintf("%vRT %v %v", a.Rt(), a.title, a.url)
				tweet(client, msg)
				time.Sleep(time.Second * 15)
			}
		}
		savecsv(articles)
	}
}

// Specify the step to count in.
func roundDown(i int) int {
	return (i / 50) * 50
}

// When refered to as Rt, the number of tweets is rounded down.
func (a Article) Rt() int {
	return roundDown(a.retweet)
}

// tweetWorthy determines if the article should be tweeted. prev is
// the number of tweets at the previous fetch.
func TweetWorthy(retweet int, prev int) bool {
	// Anything less than 100RT is not worth it, yet.
	if retweet < 100 {
		return false
	}
	r := roundDown(retweet)
	p := roundDown(prev)
	// When the retweet count surpasses 100, 150, 200, 250... and
	// so on.
	return r-p >= 50
}

func tweet(c *anaconda.TwitterApi, msg string) {
	if Config.dryrun == false {
		v := url.Values{}
		_, err := c.PostTweet(msg, v)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(msg)
}

func createMap(as []*Article) map[string]int {
	m := make(map[string]int)
	for _, a := range as {
		m[a.url] = a.retweet
	}
	return m
}

// Returns articles that have been published in the last _days_.
func withinDays(as []*Article, days int) []*Article {
	var articles []*Article

	for _, a := range as {
		d := time.Now().Sub(a.published)
		if d.Hours() <= float64(days*24) {
			articles = append(articles, a)
		}
	}
	return articles
}

func fetchArticles(pages int) []*Article {
	var articleLink = regexp.MustCompile("<a.*river_headline.*>")

	var content []byte
	var articles []*Article

	for p := 1; p <= pages; p++ {
		content = getPage(pageUrl(p))
		match := articleLink.FindAll(content, -1)

		for _, m := range match {
			a := parseArticleLink(m)
			a.retweet = tweetCount(a.url)
			articles = append(articles, a)
		}
	}
	return articles
}

// We store the retweet count and the publish date for the later reference.
func (a Article) csv() []string {
	var row []string
	var rt string

	rt = strconv.Itoa(a.retweet)
	row = append(row, a.published.Format(time.RFC3339), a.url, a.title, rt)
	return row
}

func savecsv(as []*Article) {
	file, err := os.OpenFile(datafile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	for _, a := range as {
		err := writer.Write(a.csv())
		if err != nil {
			log.Fatal(err)
		}
	}
	writer.Flush()
}

func loadcsv() []*Article {
	articles := []*Article{}

	file, err := os.Open(datafile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		n, _ := strconv.Atoi(row[3])
		t, _ := time.Parse(time.RFC3339, row[0])
		a := &Article{published: t, url: row[1], title: row[2], retweet: n}
		articles = append(articles, a)
	}
	return articles
}

// Initialize the twitter api client with Anaconda lib. Specify the
// config file.
func NewTwitterClient(filename string) *anaconda.TwitterApi {
	conf, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(conf)
	keys := map[string]string{}
	err = decoder.Decode(&keys)
	if err != nil {
		log.Fatal(err)
	}
	anaconda.SetConsumerKey(keys["consumerkey"])
	anaconda.SetConsumerSecret(keys["consumersecret"])
	c := anaconda.NewTwitterApi(keys["accesstoken"], keys["tokensecret"])
	return c
}

func pageUrl(n int) string {
	url := "http://jp.techcrunch.com"
	if n == 1 {
		return url
	} else {
		return fmt.Sprintf("%s/page/%d/", url, n)
	}
}

func getPage(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

// Create an Article object with an a tag link html fragment.
func parseArticleLink(atag []byte) *Article {
	var articleDate = regexp.MustCompile("201[0-9]/[0-9]+/[0-9]+")
	var articleURL = regexp.MustCompile("href=\"([^\"]+)\"")
	var articleTitle = regexp.MustCompile("title=\"([^\"]+)\"")

	date := articleDate.Find(atag)
	url := articleURL.FindSubmatch(atag)
	title := articleTitle.FindSubmatch(atag)
	loc, _ := time.LoadLocation("Asia/Tokyo")
	d, _ := time.ParseInLocation("2006/01/02", string(date), loc)

	return &Article{url: string(url[1]), title: string(title[1]), published: d}
}

// tweetCount returns the number of tweets that include a given url.
func tweetCount(url string) int {
	res, err := http.Get(retweetAPI + url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		log.Fatal(err)
	}
	num, _ := (result["count"]).(float64) // type assertion is necessary
	return int(num)
}

// Web UI
func handler(w http.ResponseWriter, r *http.Request) {
	d, err := ioutil.ReadFile(datafile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(d[:]))
}

func startServer() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
