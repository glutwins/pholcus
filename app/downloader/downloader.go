package downloader

import (
	"errors"
	"net/http"

	"github.com/glutwins/pholcus/app/downloader/request"
	"github.com/glutwins/pholcus/app/downloader/surfer"
	"github.com/glutwins/pholcus/config"
)

var surf = surfer.New()
var phantom = surfer.NewPhantom(config.DefaultConfig.PhantomJs, config.PHANTOMJS_TEMP)

func Download(cReq *request.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	switch cReq.GetDownloaderID() {
	case request.SURF_ID:
		resp, err = surf.Download(cReq)

	case request.PHANTOM_ID:
		resp, err = phantom.Download(cReq)
	}

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New("响应状态 " + resp.Status)
	}

	return resp, nil
}
