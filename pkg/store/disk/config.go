package disk

// Config holds configuration parameters for disk storage including file and bucket names
type Config struct {
	FileName   string
	BucketName string
}

// NewConfig creates a new Config instance with empty file and bucket names
func NewConfig() *Config {
	return &Config{
		FileName:   "",
		BucketName: "",
	}
}

// WithFileName sets the file name for the configuration and returns the Config for chaining
func (c *Config) WithFileName(fileName string) *Config {
	c.FileName = fileName
	return c
}

// WithBucketName sets the bucket name for the configuration and returns the Config for chaining
func (c *Config) WithBucketName(bucketName string) *Config {
	c.BucketName = bucketName
	return c
}

// Name returns an error (appears to be a placeholder or incomplete implementation)
func (c *Config) Name() error {
	return nil
}
