package util

import "time"

var CfgName = "cloud-cfg.json"
var CfgBase = []byte(`{
    "part_size_bytes": 2097152,
    "num_threads": 3,
    "wait_time_seconds": 5
}`)

// config settings
type Config struct {
	PartSize        int64 `json:"part_size_bytes"`
	NumThreads      int   `json:"num_threads"`
	WaitTimeSeconds int   `json:"wait_time_seconds"`
}

func (c *Config) GetPartSize() int64 {
	return c.PartSize
}
func (c *Config) GetNumThreads() int {
	return c.NumThreads
}
func (c *Config) GetWaitTime() time.Duration {
	return time.Duration(c.WaitTimeSeconds) * time.Second
}
