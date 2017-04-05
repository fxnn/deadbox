package config

// Drop configuration, created once per configured drop
type Drop struct {
	// Name identifies this drop uniquely
	Name string

	// ListenAddress defines the the local network address this drop
	// shall listen on.
	ListenAddress string
}
