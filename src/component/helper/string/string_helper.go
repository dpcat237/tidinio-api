package string_helper

import (
	"encoding/base64"
	"crypto/sha1"
	"strings"
	"net/url"
	"errors"
	"io"
	"strconv"
	"unicode/utf8"
	"github.com/microcosm-cc/bluemonday"
	"github.com/abadojack/whatlanggo"
	"github.com/djimenez/iconv-go"
	"github.com/saintfish/chardet"
)

func CleanUrl(feedUrl string) (string, error) {
	feedUrl = fixUrl(feedUrl)

	return feedUrl, validateUrl(feedUrl)
}

//Converts data from slice of bytes to UTF-8
func ConvertDataToUtf8(bytesData []byte, charset string) []byte {
	out := make([]byte, len(bytesData))
	iconv.Convert(bytesData, out, charset, "UTF-8")

	return out
}

//Detects charset from slice of bytes
func DetectDataCharset(bytesData []byte) (string, error) {
	charDetector := chardet.NewHtmlDetector()
	detected, e := charDetector.DetectBest(bytesData)
	if e != nil {
		return "", e
	}

	return detected.Charset, nil
}

func DetectLanguageFromHtml(data string) string {
	content := StripHtmlContent(data)
	if (StringLength(content) < 100) {
		return ""
	}
	info := whatlanggo.Detect(content)

	return whatlanggo.LangToString(info.Lang)
}

func GetHashFromString(text string) string {
	hasher := sha1.New()
	io.WriteString(hasher, text)

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

func StringLength(value string) int {
	return utf8.RuneCountInString(value)
}

func StringToUint(value string) uint {
	var result uint
	if bar, err := strconv.Atoi(value); err == nil {
		result = uint(bar)
	}

	return result
}

func StripHtmlContent(value string) string {
	p := bluemonday.StrictPolicy()

	return p.Sanitize(value)
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
