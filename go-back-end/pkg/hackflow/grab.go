package hackflow

import (
	"fmt"
	"os"
	"time"

	"github.com/cavaliercoder/grab"
)

type Grab struct {
	client *grab.Client
}

func NewGrab() *Grab {
	return &Grab{
		client: grab.NewClient(),
	}
}

func (g *Grab) Install(link, dstPath string) (string, error) {
	req, _ := grab.NewRequest(dstPath, link)
	// start download
	logger.Debugf("Downloading %v...\n", req.URL())
	resp := g.client.Do(req)
	logger.Debugf("  %v\n", resp.HTTPResponse.Status)
	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			logger.Debugf("\r  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}
	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}
	logger.Debugf("Download saved to %v \n", resp.Filename)
	return dstPath, nil
}
