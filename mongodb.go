package mongoout

import (
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/outputs"
)

func init() {
	outputs.RegisterType("mongodb", makeMongodb)
}

func makeMongodb(
	_ outputs.IndexManager,
	info beat.Info,
	observer outputs.Observer,
	cfg *common.Config,
) (outputs.Group, error) {

	config := defaultConfig()
	if err := cfg.Unpack(&config); err != nil {
		return outputs.Fail(err)
	}

	hosts, err := outputs.ReadHostList(cfg)
	if err != nil {
		return outputs.Fail(err)
	}

	clients := make([]outputs.NetworkClient, len(hosts))
	for i, h := range hosts {
		client, err := newClient(h, config.DB, config.Collection, observer, info, config.Timeout)
		if err != nil {
			return outputs.Fail(err)
		}
		clients[i] = newBackoffClient(client, config.Backoff.Init, config.Backoff.Max)
	}

	return outputs.SuccessNet(config.LoadBalance, config.bulkMaxSize, config.MaxRetries, clients)
}
