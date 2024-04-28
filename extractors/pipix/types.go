package pipix

type pipixData struct {
	Data map[string]struct {
		ItemId    int64  `json:"item_id"`
		ItemIdStr string `json:"item_id_str"`
		Content   string `json:"content"`
		Comments  []struct {
			Text string `json:"text"`
			Item struct {
				Content string `json:"content"`
				Video   struct {
					VideoHigh struct {
						Width   int `json:"width"`
						Height  int `json:"height"`
						UrlList []struct {
							Url     string `json:"url"`
							Expires int64  `json:"expires"`
						} `json:"url_list"`
						CoverImage struct {
							Width   int `json:"width"`
							Height  int `json:"height"`
							UrlList []struct {
								Url string `json:"url"`
							} `json:"url_list"`
						}
					} `json:"video_high"`
				} `json:"video"`
			} `json:"item"`
		} `json:"comments"`
		Video struct {
			VideoWidth    int `json:"video_width"`
			VideoHeight   int `json:"video_height"`
			VideoDownload struct {
				Width   int `json:"width"`
				Height  int `json:"height"`
				UrlList []struct {
					Url     string `json:"url"`
					Expires int64  `json:"expires"`
				} `json:"url_list"`
			} `json:"video_download"`
		} `json:"video"`
	}
}
