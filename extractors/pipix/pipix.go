package pipix

import (
	_ "embed"
	"encoding/json"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/pkg/errors"
	"net/http"
	"strconv"

	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/request"
	"github.com/iawia002/lux/utils"
)

func init() {
	e := New()
	extractors.Register("pipix", e)
}

type extractor struct{}

// New returns a douyin extractor.
func New() extractors.Extractor {
	return &extractor{}
}

// Extract is the main function to extract the data.
func (e *extractor) Extract(url string, option extractors.Options) ([]*extractors.Data, error) {
	fmt.Println("pipix start")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close() // nolint
	url = resp.Header.Get("location")

	fmt.Println("loca Url", url)

	itemIds := utils.MatchOneOf(url, `/item/(\d+)`)
	if len(itemIds) == 0 {
		return nil, errors.New("unable to get video ID")
	}
	if itemIds == nil || len(itemIds) < 2 {
		return nil, errors.WithStack(extractors.ErrURLParseFailed)
	}
	itemId := itemIds[len(itemIds)-1]

	// dynamic generate cookie
	//cookie, err := createCookie()
	//if err != nil {
	//	return nil, errors.WithStack(err)
	//}

	api := "https://h5.pipix.com/bds/webapi/item/detail/?item_id=" + itemId + "&source=share"

	// define request headers and sign agent
	headers := map[string]string{}
	//headers["Cookie"] = cookie
	headers["Referer"] = "https://h5.pipix.com/"
	headers["User-Agent"] = browser.Chrome()

	jsonData, err := request.Get(api, url, headers)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var pipix pipixData
	if err = json.Unmarshal([]byte(jsonData), &pipix); err != nil {
		return nil, errors.WithStack(err)
	}

	fmt.Println("pipix", pipix)

	mainData, ok := pipix.Data["item"]
	if !ok {
		return nil, errors.New("unable to get video Data")
	}

	title := mainData.Content
	streams := make(map[string]*extractors.Stream)

	if len(mainData.Comments) == 0 {
		videoList := mainData.Video.VideoDownload.UrlList
		for i, v := range videoList {
			totalSize, _ := request.Size(v.Url, url)

			streams[strconv.Itoa(i)] = &extractors.Stream{
				Quality: fmt.Sprintf("%d*%d", mainData.Video.VideoDownload.Width, mainData.Video.VideoDownload.Height),
				Parts: []*extractors.Part{
					{
						URL:  v.Url,
						Size: totalSize,
						Ext:  "mp4",
					},
				},
			}
		}
	} else {
		videoList := mainData.Comments[0].Item.Video.VideoHigh.UrlList
		for i, v := range videoList {
			totalSize, _ := request.Size(v.Url, url)

			streams[strconv.Itoa(i)] = &extractors.Stream{
				Quality: fmt.Sprintf("%d*%d", mainData.Comments[0].Item.Video.VideoHigh.Width, mainData.Comments[0].Item.Video.VideoHigh.Height),
				Parts: []*extractors.Part{
					{
						URL:  v.Url,
						Size: totalSize,
						Ext:  "mp4",
					},
				},
			}
		}
	}

	return []*extractors.Data{
		{
			Site:    "皮皮虾 pipix.com",
			Title:   title,
			Type:    "video",
			Streams: streams,
			URL:     url,
		},
	}, nil
}
