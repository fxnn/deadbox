package config

import "net/url"

func Dummy() *Application {
	commanderUrl, _ := url.Parse("http://localhost:6545")
	a := Agent{
		CommanderUrls: []*url.URL{commanderUrl},
	}
	c := Commander{
		ListenAddress: ":6545",
	}
	app := &Application{
		Agents:     []Agent{a},
		Commanders: []Commander{c},
	}
	return app
}
