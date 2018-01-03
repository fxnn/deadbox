package application

import (
	"path/filepath"

	"github.com/fxnn/deadbox/config"
)

const (
	filePermUserAndGroupCanReadOrWrite = 0660
	filePermOnlyUserCanReadOrWrite     = 0600
	dbFileExtension                    = "boltdb"
	privateKeyFileExtension            = "pem"
	certFileExtension                  = "pem"
)

func dbFileName(cfg *config.Application, name string) string {
	return filepath.Join(cfg.DbPath, name+"."+dbFileExtension)
}

func privateKeyFileName(path string, name string) string {
	return filepath.Join(path, name+".private."+privateKeyFileExtension)
}

func certFileName(path string, dropName string) string {
	return filepath.Join(path, dropName+"."+certFileExtension)
}
