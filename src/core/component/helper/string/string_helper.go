package string_helper

import (
	"net/http"
	"io/ioutil"
	"encoding/base64"
	"crypto/sha1"
	"strings"
	"net/url"
	"errors"
	"io"
)

func CleanUrl(feedUrl string) (string, error) {
	feedUrl = fixUrl(feedUrl)

	return feedUrl, validateUrl(feedUrl)
}

func GetHashFromString(text string) string {
	hasher := sha1.New()
	io.WriteString(hasher, text)

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func GetHashFromUrlData(url string) string {
	resp, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(resp.Body)
	hasher := sha1.New()
	hasher.Write(bytes)

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func fixUrl(feedUrl string) string {
	if (!strings.Contains(feedUrl, "://")) {
		feedUrl = "http://" + feedUrl
	}
	feedUrl = strings.Replace(feedUrl, "feed:", "http:", 1)
	if last := len(feedUrl) - 1; last >= 0 && feedUrl[last] == '/' {
		feedUrl = feedUrl[:last]
	}

	return feedUrl
}

func validateUrl(feedUrl string) error {
	base, err := url.Parse(feedUrl)
	if err != nil {
		return err
	}

	if (base.Scheme == "http" || base.Scheme == "feed" || base.Scheme == "https") {
		return nil
	}

	return errors.New("Wrong url")
}
