package config

// Application configuration, created once in application lifecycle.
type Application struct {

	// Agents configured in this application
	Agents []Agent

	// Commanders configured in this application
	Commanders []Commander

}
