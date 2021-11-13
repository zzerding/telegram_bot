package bingwallpaper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

// get cn.bing.com json strct
type BingResponse struct {
	Images   interface{} `json:"images"`
	Tooltips interface{} `json:"tooltips"`
}

//set error type
var (
	// response json parer error
	errJsonParse error = errors.New("json parser error")

	// request server error
	errBingServer error = errors.New("connection bing server error")

	//cache expired
	errCacheExpired error = errors.New("cache key is expired")
	UrlBingServer         = "https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1"
)

//cache bing request everyday
var CacheBingRequest = cache.New(5*time.Minute, 10*time.Minute)

//cache key
const cacheKEY string = "cacheBingRequestKEY"

//get cache value
func GetCache(url string) (interface{}, error) {
	value, found := CacheBingRequest.Get(cacheKEY)
	if found {
		return value, nil
	}
	data, err := getBingWallPaper(url)
	if err != nil {
		return nil, errCacheExpired
	}
	CacheBingRequest.Add(cacheKEY, data, 1*time.Hour)
	return data, nil
}

//delete cache
func deleteCache() {
	CacheBingRequest.Delete(cacheKEY)
}

//http clint
func request(method string, url string, body io.Reader) (*http.Response, error) {
	ck := &http.Cookie{
		Name:  "ENSEARCH",
		Value: "BENVER=0",
	}
	clint := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic("create request error")
	}
	req.AddCookie(ck)
	req.Header.Add("uhd", "1")
	req.Header.Add("uhdwidth", "3840")
	req.Header.Add("uhdheight", "2160")
	req.Header.Add("nc", fmt.Sprint(time.Now().Unix()*1000))
	resp, err := clint.Do(req)
	return resp, err
}

//get bingwallpaper
func getBingWallPaper(url string) (BingResponse, error) {
	result := &BingResponse{}
	// r, err := http.Get(url)
	r, err := request(http.MethodGet, url, nil)
	defer r.Body.Close()
	if err != nil || r.StatusCode != http.StatusOK {
		return *result, errBingServer
	}
	dec := json.NewDecoder(r.Body)
	decodeErr := dec.Decode(result)
	if decodeErr != nil {
		return *result, errJsonParse
	}
	return *result, nil
}
