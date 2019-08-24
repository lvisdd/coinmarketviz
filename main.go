package main

import (
	"net/http"
    "time"
	"encoding/json"
    "io/ioutil"
	gocache "github.com/patrickmn/go-cache"
	"github.com/unrolled/render"
	"google.golang.org/appengine"
	// "google.golang.org/appengine/log"
    "google.golang.org/appengine/urlfetch"
)

var (
	cache = gocache.New(1*time.Minute, 30*time.Second)
)

type Ticker struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	Rank             string `json:"rank"`
	PriceUsd         string `json:"price_usd"`
	PriceBtc         string `json:"price_btc"`
	Two4HVolumeUsd   string `json:"24h_volume_usd"`
	MarketCapUsd     string `json:"market_cap_usd"`
	AvailableSupply  string `json:"available_supply"`
	TotalSupply      string `json:"total_supply"`
	MaxSupply        string `json:"max_supply"`
	PercentChange1H  string `json:"percent_change_1h"`
	PercentChange24H string `json:"percent_change_24h"`
	PercentChange7D  string `json:"percent_change_7d"`
	LastUpdated      string `json:"last_updated"`
}

type Global struct {
    total_market_cap_usd             string `json:"total_market_cap_usd"` 
    total_24h_volume_usd             string `json:"total_24h_volume_usd"`
    bitcoin_percentage_of_market_cap string `json:"bitcoin_percentage_of_market_cap"`
    active_currencies                string `json:"active_currencies"`
    active_assets                    string `json:"active_assets"`
    active_markets                   string `json:"active_markets"`
    last_updated                     string `json:"last_updated"`
}

// func getJson(this interface{}, url string, r *http.Request) error {
func getJson(this interface{}, url string, r *http.Request) (interface{}, error) {
    ctx := appengine.NewContext(r)
    client := urlfetch.Client(ctx)
    res, err := client.Get(url)
	// res, err := http.Get(url)
	// log.Debugf(ctx, "%+v", err)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

    var ticker []Ticker
  	err = json.Unmarshal(b, &ticker)
	if err != nil {
		return nil, err
	}

	// return json.NewDecoder(res.Body).Decode(this)

	// allCoins := make(map[string]Ticker)
	// for i := 0; i < len(ticker); i++ {
	// 	allCoins[ticker[i].ID] = ticker[i]
    // 	log.Debugf(ctx, "%+v", ticker[i])
	// }
	return ticker, nil
}

// func loadData(this *Data, url string, r *http.Request) (interface{}, bool) {
func loadData(this interface{}, url string, r *http.Request) (interface{}, bool) {
    // ctx := appengine.NewContext(r)
	if cached, found := cache.Get("data"); found {
		cache.Set("data_tmp", cached, gocache.NoExpiration)
		return cached, found
	}

	if cached_tmp, found_tmp := cache.Get("data_tmp"); found_tmp {
		go func() {
			// getJson(this, url)
			ticker, _ := getJson(this, url, r)
			// cache.Set("data", this, gocache.DefaultExpiration)
			cache.Set("data", ticker, gocache.DefaultExpiration)
		}()
		return cached_tmp, found_tmp
	}

	// getJson(this, url)
	ticker, _ := getJson(this, url, r)

	// cache.Set("data", this, gocache.DefaultExpiration)
	cache.Set("data", ticker, gocache.DefaultExpiration)

	// return this, false
	return ticker, false
}

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Vary", "Accept-Encoding")
}

func setCacheHeader(w http.ResponseWriter, found bool) {
	v := "MISS"
	if found {
		v = "HIT"
	}
	w.Header().Set("X-Cache", v)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ticker/", tickerHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}

func tickerHandler(w http.ResponseWriter, r *http.Request) {
	render := render.New()
    // ctx := appengine.NewContext(r)

    setDefaultHeaders(w)

	// data := new(Data)
	ticker := []Ticker{}
	url := "https://api.coinmarketcap.com/v1/ticker/"

	// res, found := loadData(data, url)
	res, found := loadData(ticker, url, r)
	setCacheHeader(w, found)

    // log.Debugf(ctx, "%+v", res)
    // log.Debugf(ctx, "%v", found)
	render.JSON(w, http.StatusOK, res)
}
