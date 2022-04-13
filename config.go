package mongoout

import "time"

type config struct {
	Index       string        `config:"index"`
	LoadBalance bool          `config:"loadbalance"`
	MaxRetries  int           `config:"max_retries"`
	Timeout     time.Duration `config:"timeout"`
	Hosts       []string      `config:"urls"`
	DB          string        `config:"db"`
	Collection  string        `config:"collection"`
	bulkMaxSize int           `config:"bulk_max_size"`
	Backoff     backoff       `config:"backoff"`
}

type backoff struct {
	Init time.Duration
	Max  time.Duration
}

func defaultConfig() config {
	return config{
		Hosts:       []string{"mongodb://localhost:27017"},
		Timeout:     5 * time.Second,
		LoadBalance: true,
		DB:          "test",
		Collection:  "test",
		bulkMaxSize: 100,
		MaxRetries:  3,
		Backoff: backoff{
			Init: 1 * time.Second,
			Max:  60 * time.Second,
		},
	}
}
