package http_helper

import (
	"bytes"
	"crypto/sha1"
	"io"
	"io/ioutil"
	"net/http"
	"encoding/base64"
	"github.com/dpcat237/articletext"
	"github.com/tidinio/src/core/component/helper/string"
)

func GetContentFromUrl(url string) string {
	reader, err := getReaderFromUrl(url)
	if err != nil {
		return ""
	}

	htmlText, err := articletext.GetArticleHtmlFromReader(reader)
	if err != nil {
		return ""
	}

	return htmlText
}

func GetHashFromUrlData(url string) string {
	resp, _ := http.Get(url)
	bytesData, _ := ioutil.ReadAll(resp.Body)
	hasher := sha1.New()
	hasher.Write(bytesData)

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func getReaderFromUrl(url string) (io.Reader, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	bytesData, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return nil, e
	}
	charset, e := string_helper.DetectDataCharset(bytesData)
	if e != nil {
		return nil, e
	}

	if (charset == "UTF-8") {
		return bytes.NewReader(bytesData), nil
	}

	return bytes.NewReader(string_helper.ConvertDataToUtf8(bytesData, charset)), nil
}
