# beats-output-mongodb


Elastic Beats Output to MongoDB.

# Usage

``` shell
git clone https://github.com/huawen0327/beats-output-mongodb.git
cd beats-output-mongodb/cmd 
go build -o filebeat.exe
./filebeat -c config.yml -e
```

### or

``` go
package main

import (
	"github.com/elastic/beats/v7/filebeat/cmd"
	inputs "github.com/elastic/beats/v7/filebeat/input/default-inputs"
	_ "github.com/huawen0327/beats-output-mongodb"
	"os"
)

func main() {
	if err := cmd.Filebeat(inputs.Init, cmd.FilebeatSettings()).Execute(); err != nil {
		os.Exit(1)
	}
}
```

### config.yml

``` yaml
output.mongodb:
  hosts: ["mongodb://localhost:27017", "mongodb://localhost:27018"]
  db: test_log
  collection: test_log
  bulk_max_size: 20
```
