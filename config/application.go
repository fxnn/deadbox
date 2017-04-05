package config

// DefaultPort is the port used by the application per default.
// It was selected such that it doesn't seem that other applications use it as
// a default.
const DefaultPort = "6545"

// Application configuration, created once in application lifecycle.
type Application struct {
	// DbFile points to a file were the database is stored
	DbFile string

	// Workers configured in this application
	Workers []Worker

	// Drops configured in this application
	Drops []Drop
}
