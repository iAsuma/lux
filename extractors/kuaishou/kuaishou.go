package kuaishou

import (
	"bytes"
	"encoding/json"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/iawia002/lux/utils"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/request"
)

func init() {
	extractors.Register("kuaishou", New())
}

type extractor struct{}

// New returns a kuaishou extractor.
func New() extractors.Extractor {
	return &extractor{}
}

// fetch url and get the cookie that write by server
func fetchCookies(url string, headers map[string]string) (string, error) {
	res, err := request.Request(http.MethodGet, url, nil, headers)
	if err != nil {
		return "", err
	}

	defer res.Body.Close() // nolint

	cookiesArr := make([]string, 0)
	cookies := res.Cookies()

	for _, c := range cookies {
		cookiesArr = append(cookiesArr, c.Name+"="+c.Value)
	}

	return strings.Join(cookiesArr, "; "), nil
}

// Extract is the main function to extract the data.
func (e *extractor) Extract(url string, option extractors.Options) ([]*extractors.Data, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url2 := url
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close() // nolint
	url = resp.Header.Get("location")

	if strings.Contains(url, "v.douyin.com") {

	} else {
		headers := map[string]string{
			"User-Agent": browser.Computer(),
			"Referer":    url2,
			"Host":       "www.kuaishou.com",
			"Cookie":     "did=web_",
		}

		dataString, err := request.Get(url, url, headers)
		utils.WriteFile("zzz.html", dataString)
		fmt.Println("ddddd", err, dataString)
		return nil, nil
	}

	pUrl, err := neturl.ParseRequestURI(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	paramsMap := make(map[string]string)
	for key, values := range pUrl.Query() {
		paramsMap[key] = values[0]
	}
	paramsMap["h5Domain"] = "v.m.chenzhongtech.com"
	paramsMap["shareChannel"] = "share_copylink"
	paramsMap["env"] = "SHARE_VIEWER_ENV_TX_TRICK"

	jsonData, err := json.Marshal(paramsMap)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	apiUrl := "https://v.m.chenzhongtech.com/rest/wd/photo/info"
	headers := map[string]string{
		"User-Agent":   browser.Mobile(),
		"Referer":      url,
		"Content-Type": "application/json; charset=UTF-8",
		"Cookie":       "did=web_",
	}

	resBody, err := request.PostByte(apiUrl, bytes.NewBuffer(jsonData), headers)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var ddDATA map[string]interface{}
	err = json.Unmarshal(resBody, &ddDATA)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var data kuiashouApiData
	if err = json.Unmarshal(resBody, &data); err != nil {
		return nil, errors.WithStack(err)
	}

	streams := make(map[string]*extractors.Stream, len(data.Photo.MainMvUrls))
	for i, v := range data.Photo.MainMvUrls {
		size, err := request.Size(v.Url, apiUrl)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		urlData := &extractors.Part{
			URL:  v.Url,
			Size: size,
			Ext:  "mp4",
		}
		streams[strconv.Itoa(i)] = &extractors.Stream{
			Parts:   []*extractors.Part{urlData},
			Size:    size,
			Quality: strconv.Itoa(i),
		}
	}

	return []*extractors.Data{
		{
			Site:    "快手 kuaishou.com",
			Title:   data.Photo.Caption,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     url,
		},
	}, nil
}
