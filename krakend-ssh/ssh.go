package ssh

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"gopkg.in/yaml.v2"
)

const Namespace = "github.com/jozefiel/krakend-ssh"
const logPrefix = "[SERVICE: ssh]"

type ssh_clients struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// New creates a new metrics producer
func New(ctx context.Context, extraConfig config.ExtraConfig, logger logging.Logger) error {

	logger.Debug(logPrefix, "Parsing ssh config")

	// var cfg *Config
	cfg, ok := configGetter(extraConfig).(sshConfig)
	if !ok {
		logger.Error("ssh has no config file:", ok)
	}
	logger.Debug(logPrefix, "hahaha", cfg)

	config_file, err := ioutil.ReadFile(cfg.config_path)
	if err != nil {
		logger.Error("ssh has no config file:", err.Error())
	}

	data := make(map[string]ssh_clients)
	err = yaml.Unmarshal(config_file, &data)
	if err != nil {
		logger.Error(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/{uniconfig}", sendCommand).Methods("PUT")
	router.HandleFunc("/api/{uniconfig}", getOutput).Methods("GET")
	// router.PathPrefix("/").Handler(http.FileServer(http.Dir("./front/")))

	for client, config := range data {
		config = parseConfigs(client, config)
		handler := &sshHandler{host: client, addr: config.Host, port: config.Port, user: config.User, secret: config.Password}
		handler.sshClient()
	}

	log.Fatal(http.ListenAndServe(":"+cfg.port, router))

	return nil
}

func parseConfigs(client_name string, config ssh_clients) ssh_clients {

	if config.Host == "" {
		log.Fatal("host can not be empty")
	}

	if config.Port == "" {
		log.Println("Port set to 22 for client", client_name)
		config.Port = "22"
	}

	if config.User == "" {
		log.Fatal("user can not be empty")
	}

	if config.Password == "" {
		log.Fatal("password can not be empty")
	}
	return config
}
