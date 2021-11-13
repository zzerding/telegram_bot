package bingwallpaper

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/patrickmn/go-cache"
)

var testData = BingResponse{Images: "1", Tooltips: "1"}
var jsonData, _ = json.Marshal(testData)

func TestBingWallPaper(t *testing.T) {

	t.Run("test get json", func(t *testing.T) {
		// t.SkipNow()
		t.Log("get json")
		url := "https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1"
		r, _ := request(http.MethodGet, url, nil)
		defer r.Body.Close()
		dec := json.NewDecoder(r.Body)
		type BingResponse struct {
			Images   interface{} `json:"images"`
			Tooltips interface{} `json:"tooltips"`
		}
		v := &BingResponse{}
		err := dec.Decode(v)
		if err != nil {
			t.Log("decode json error")
			t.Log(err)
			t.Errorf("get json error %v", err)
		}
		t.Log("data")
		t.Log(v)
	})

	t.Run("test getBingWallPaper url", func(t *testing.T) {
		got := testData
		server := mokeServer(jsonData, http.StatusOK)
		defer server.Close()
		want, err := getBingWallPaper(server.URL)
		if err != nil {
			t.Errorf("getBingWallPaper want err is nil %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf(`got %v want %v`, got, want)
		}
	})
	t.Run("test getBingWallPaper conection error", func(t *testing.T) {
		s := mokeServer(jsonData, http.StatusAccepted)
		_, err := getBingWallPaper(s.URL)
		if err != errBingServer {
			t.Errorf(`got "%s" err want "%s"`, err, errBingServer)
		}
	})
	t.Run("test getBingWallPaper json parse error", func(t *testing.T) {
		s := mokeServer([]byte("fdafsf"), http.StatusOK)
		_, err := getBingWallPaper(s.URL)
		if err != errJsonParse {
			t.Errorf(`got "%s" err want "%s"`, err, errJsonParse)
		}
	})

}

func TestGetCache(t *testing.T) {
	s := mokeServer(jsonData, http.StatusOK)
	deleteCache()
	t.Run("test getCache noCache", func(t *testing.T) {
		got, err := GetCache(s.URL)
		want := testData
		if err != nil {
			t.Error("error want nil")
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf(`got %T"%s" want %T"%s" `, got, got, want, want)
		}
	})
	t.Run("test getCache have cache", func(t *testing.T) {
		CacheBingRequest.Set(cacheKEY, jsonData, cache.DefaultExpiration)
		got, _ := GetCache(s.URL)
		want := testData
		if !reflect.DeepEqual(got, jsonData) {
			t.Errorf(`got "%s" want "%s" `, got, want)
		}
	})
	t.Run("test getCache error", func(t *testing.T) {
		s := mokeServer(jsonData, http.StatusAccepted)
		deleteCache()
		_, got := GetCache(s.URL)
		if got != errCacheExpired {
			t.Error("getCache want cache expired")
		}
	})
}

//moke server
func mokeServer(data []byte, status int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(data)
	}))
	return server
}
