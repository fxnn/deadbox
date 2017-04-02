package config

import "net/url"

func Dummy() *Application {
	dropUrl, _ := url.Parse("http://localhost:6545")
	w := Worker{
		DropUrls: []*url.URL{dropUrl},
	}
	d := Drop{
		ListenAddress: ":6545",
	}
	app := &Application{
		Workers: []Worker{w},
		Drops:   []Drop{d},
	}
	return app
}
