package config

// Application configuration, created once in application lifecycle.
type Application struct {

	// Workers configured in this application
	Workers []Worker

	// Drops configured in this application
	Drops []Drop

}
