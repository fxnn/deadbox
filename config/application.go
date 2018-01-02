package config

// DefaultPort is the port used by the application per default.
// It was selected such that it doesn't seem that other applications use it as
// a default.
const DefaultPort = "6545"

// Application configuration, created once in application lifecycle.
type Application struct {
	// DbPath points to a directory were the database files are stored in.
	// The files are named after the Drop / Worker name and created if not existing.
	DbPath string

	// PrivateKeyPath refers to a path containing private key files for each Drop and Worker.
	// The files will contain ASN.1 PKCS#1 DER encoded private RSA keys.
	// They are named after the Drop / Worker name and created if not existing.
	PrivateKeyPath string

	// Workers configured in this application
	Workers []Worker

	// Drops configured in this application
	Drops []Drop
}
