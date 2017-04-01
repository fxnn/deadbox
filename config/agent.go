package config

import "net/url"

// Agent configuration, created once per configured agent.
type Agent struct {

	// Name identifies this agent uniquely
	Name string

	// CommanderUrls identifies the commander instances this agent should connect to
	CommanderUrls []*url.URL

}
