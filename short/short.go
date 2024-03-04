package short

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Shorter struct {
}

func (s Shorter) GetLink(urlUploaded string) (string, error) {
	values := url.Values{
		"url": {urlUploaded},
	}
	shorterService := os.Getenv("SHORTEN_SERVICE")
	if len(shorterService) == 0 {
		shorterService = "https://s.inxo.ru/shorten"
	}
	response, err := http.PostForm(shorterService, values)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(response.Body)
	var js ShortUrl
	body, err := io.ReadAll(response.Body)
	err = json.Unmarshal(body, &js)
	if err != nil {
		return "", err
	}
	return js.Url, nil
}

type ShortUrl struct {
	Origin string `json:"origin"`
	Url    string `json:"short"`
}

func NewLink(uploaded string) (string, error) {
	s := Shorter{}
	return s.GetLink(uploaded)
}
