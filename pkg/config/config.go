package config

import "time"

type RPCConfig struct {
	Timeout              time.Duration `yaml:"timeout"`
	Interval             time.Duration `yaml:"interval"`
	TooManyRequestsDelay time.Duration `yaml:"too_many_requests_delay"`
}

type StorageConfig struct {
	Retention            time.Duration `yaml:"retention"`
	CleanInterval        time.Duration `yaml:"clean_interval"`
	StoreAllTransactions bool          `yaml:"store_all_transactions"`
}

type Config struct {
	Storage *StorageConfig `yaml:"storage"`
	RPC     *RPCConfig     `yaml:"rpc"`
}
