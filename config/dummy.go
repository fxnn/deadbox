package config

import "net/url"

func Dummy() *Application {
	dropUrl, _ := url.Parse("http://localhost:" + DefaultPort)
	w := Worker{
		Name:     "Default Worker",
		DropUrls: []*url.URL{dropUrl},
	}
	d := Drop{
		Name:          "Default Drop",
		ListenAddress: ":" + DefaultPort,
	}
	app := &Application{
		DbFile:  "deadbox.boltdb",
		Workers: []Worker{w},
		Drops:   []Drop{d},
	}
	return app
}
