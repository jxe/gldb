package gldb

import (
	"encoding/json"
	"net/http"
	"net/url"
	"io/ioutil"
)

func jsonFromURL(u string, v interface{}) (err error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return
	}
	resp, err := http.Get(u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if len(parsed.Fragment) > 0 {
		var f interface{}
    	err = json.Unmarshal(body, &f)
    	if err != nil {
    		return
    	}
    	m := f.(map[string]interface{})
    	body, err = json.Marshal(m[parsed.Fragment])
    	if err != nil {
    		return
    	}
	}
	err = json.Unmarshal(body, v)
	return
}
