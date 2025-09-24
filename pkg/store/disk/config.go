package disk

import (
	"github.com/happyhackingspace/vulnerable-target/pkg/store/config"
)

// Config holds configuration parameters for disk storage including file and bucket names
type Config struct {
	FileName   string
	BucketName string
	config.Struct
}
