package ssh

import (
	"errors"

	"github.com/luraproject/lura/v2/config"
)

type sshConfig struct {
	config_path string
	port        string
	api         string
}

func configGetter(extraConfig config.ExtraConfig) interface{} {
	value, ok := extraConfig[Namespace]
	if !ok {
		return nil
	}

	castedConfig, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}

	cfg := sshConfig{}

	if value, ok := castedConfig["config_path"]; ok {
		cfg.config_path = value.(string)
	}

	if value, ok := castedConfig["port"]; ok {
		cfg.port = value.(string)
	}

	if value, ok := castedConfig["api"]; ok {
		cfg.api = value.(string)
	} else {
		cfg.api = "api"
	}

	return cfg
}

var ErrNoConfig = errors.New("ssh: unable to load custom config")
