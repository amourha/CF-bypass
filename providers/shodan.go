package providers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/reusee/mmh3"
)

const (
	ShodanAPIURL = "https://api.shodan.io/shodan/host/search?key=%s&query=http.favicon.hash:%d&facets=asn,city,country,domain,isp,org,port,state"
	FaviconURL   = "http://%s/favicon.ico"
)

type Shodan struct {
	config *Config
}

//Base64Split splits a base64 string according to RFC 2045
func Base64Split(s string) string {
	ss := ""
	size := 76
	for len(s) > 0 {
		if len(s) < size {
			size = len(s)
		}
		ss += s[:size] + "\n"
		s = s[size:]

	}
	return ss
}

func (s *Shodan) GetFavicon(domain string) ([]byte, error) {
	RequestURL := fmt.Sprintf(FaviconURL, domain)
	resp, err := s.config.Client.DoRequest(RequestURL)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return body, err
}

//CalcFaviconHash calculates Favicon Hash
func CalcFaviconHash(favicon []byte) int32 {
	var base64Favicon = Base64Split(base64.StdEncoding.EncodeToString(favicon))
	return int32(mmh3.Hash32([]byte(base64Favicon)))
}

//MakeShodanAPICall Makes an API call to search for the favicon
func (s *Shodan) MakeShodanAPICall(hash int32, APIKey string) ([]byte, error) {
	RequestURL := fmt.Sprintf(ShodanAPIURL, APIKey, hash)
	resp, err := s.config.Client.DoRequest(RequestURL)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return body, err
}

func NewShodan(c *Config) Provider {
	return &Shodan{config: c}
}

func (s *Shodan) BypassCF(domain string, results chan<- string) error {
	favicon, err := s.GetFavicon(domain)
	if err != nil {
		return err
	}
	faviconHash := CalcFaviconHash(favicon)
	result, err := s.MakeShodanAPICall(faviconHash, s.config.ShodanAPIKey)
	if err != nil {
		return err
	}

	results <- string(result)
	return nil
}
