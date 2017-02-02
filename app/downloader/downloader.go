package downloader

import (
	"errors"
	"net/http"

	"github.com/glutwins/pholcus/app/downloader/request"
	"github.com/glutwins/pholcus/app/downloader/surfer"
	"github.com/glutwins/pholcus/app/spider"
	"github.com/glutwins/pholcus/config"
)

// The Downloader interface.
// You can implement the interface by implement function Download.
// Function Download need to return Page instance pointer that has request result downloaded from Request.
type Downloader interface {
	Download(*spider.Spider, *request.Request) (*spider.Context, error)
}

type Surfer struct {
	surf    surfer.Surfer
	phantom surfer.Surfer
}

var SurferDownloader = &Surfer{
	surf:    surfer.New(),
	phantom: surfer.NewPhantom(config.DefaultConfig.PhantomJs, config.PHANTOMJS_TEMP),
}

func (self *Surfer) Download(sp *spider.Spider, cReq *request.Request) (*spider.Context, error) {
	var resp *http.Response
	var err error

	switch cReq.GetDownloaderID() {
	case request.SURF_ID:
		resp, err = self.surf.Download(cReq)

	case request.PHANTOM_ID:
		resp, err = self.phantom.Download(cReq)
	}

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New("响应状态 " + resp.Status)
	}

	ctx := spider.GetContext(sp, cReq)
	ctx.Response = resp

	return ctx, nil
}
