package short

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"io"
	"net/http"
	"net/url"
)

type Shorter struct {
	myApp fyne.App
}

func (s Shorter) GetLink(urlUploaded string) (string, error) {
	values := url.Values{
		"url": {urlUploaded},
	}
	shorterService := s.myApp.Preferences().String("SHORT_SERVICE")
	if len(shorterService) == 0 {
		shorterService = "https://s.n32b.com/shorten"
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

func NewLink(uploaded string, app fyne.App) (string, error) {
	s := Shorter{myApp: app}
	return s.GetLink(uploaded)
}
