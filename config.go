package bartender

import (
     "encoding/json"
     "os"
     "fmt"
)

type config struct {
	DomainName string
     Port string
     SiteName string
     DBUser string
     DBName string
     DBPassword string
}


func loadConfig(path string) (*config, error) {
	configFile, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("Unable to read configuration file %s", path)
	}

	config := new(config)

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse configuration file %s", path)
	}

	return config, nil
}
