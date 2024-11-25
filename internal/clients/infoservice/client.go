package infoservice

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"online-song-library/internal/config"
	msong "online-song-library/internal/model/song"

	"github.com/sirupsen/logrus"
)

type Client struct {
	ReqURL     *url.URL
	HTTPClient *http.Client
}

func NewMusicInfoClient(cfg *config.Config) *Client {
	reqURL, err := url.Parse(cfg.MusicInfoURL)
	if err != nil {
		logrus.Errorf("error parsing URL: %v\n", err)
	}

	return &Client{
		ReqURL: reqURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetSongInfo(song msong.Song) (map[string]string, error) {
	var response struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}

	query := c.ReqURL.Query()
	query.Set("group", song.Group)
	query.Set("song", song.Song)
	c.ReqURL.RawQuery = query.Encode()

	resp, err := c.HTTPClient.Get(c.ReqURL.String())
	if err != nil {
		logrus.Errorf("error making GET request: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("unexpected status code %d\n", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		logrus.Errorf("error decoding JSON response: %v\n", err)
	}

	return map[string]string{
		"ReleaseDate": response.ReleaseDate,
		"Text":        response.Text,
		"Link":        response.Link,
	}, nil
}
