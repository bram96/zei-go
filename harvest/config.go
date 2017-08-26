package harvest

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func readConfig() (Config, error) {
	var config Config
	c, err := ioutil.ReadFile("config/harvest.yml")
	if err != nil {
		return config, fmt.Errorf("Could not open harvest-config.yaml: %v", err)
	}
	if err := yaml.Unmarshal(c, &config); err != nil {
		return config, fmt.Errorf("Could not parse harvest-config.yaml: %v", err)
	}
	return config, nil
}

func writeConfig(c Config) error {
	config, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	if err = ioutil.WriteFile("config/harvest.yml", config, 0644); err != nil {
		return fmt.Errorf("Failed to write config to file: %v", err)
	}
	return nil
}
