package config

import "net/url"

func Dummy() *Application {
	dropUrl, _ := url.Parse("http://localhost:" + DefaultPort)
	w := Worker{
		Name:    "Default Worker",
		DropUrl: dropUrl,
	}
	d := Drop{
		Name:          "Default Drop",
		ListenAddress: ":" + DefaultPort,
	}
	app := &Application{
		DbPath:  "./",
		Workers: []Worker{w},
		Drops:   []Drop{d},
	}
	return app
}
