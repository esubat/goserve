package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port       int    `yaml:"port"`
	Target     string `yaml:"target"`
	StaticPath string `yaml:"static_path"`
}

func main() {
	serveMode := flag.String("serve", "static", "server mode: 'proxy' or 'static'")
	configPath := flag.String("config", "config.yaml", "path to configuration file")

	flag.Parse()

	serverConfig, err := parseConfig(*configPath)

	if err != nil {
		log.Printf("Failed to parse config, using defaults: %v", err)
		serverConfig = Config{Port: 8080, Target: "http://localhost:3000", StaticPath: "./public"}
	}

	switch *serveMode {
	case "proxy":
		http.Handle("/", handleProxy(serverConfig))
	case "static":
		fs := http.FileServer(http.Dir(serverConfig.StaticPath))
		http.Handle("/", fs)
	default:
		log.Fatalf("Unknown serve type: %s. Valid options are 'proxy' or 'static-server'", *serveMode)

	}

	log.Printf("starting server at %v\n", serverConfig.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", serverConfig.Port), nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

func handleProxy(option Config) http.Handler {
	targetURL, _ := url.Parse(option.Target)
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ModifyResponse = modifyResponse()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s path: %s Method: %s", time.Now(), r.URL.Path, r.Method)
		proxy.ServeHTTP(w, r)
	})

}

func modifyResponse() func(*http.Response) error {
	return func(res *http.Response) error {
		res.Header.Set("X-Proxy", "Goserve")
		res.Header.Set("Server", "goserve")
		return nil
	}
}

func parseConfig(path string) (Config, error) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("error reading YAML file: %v", err)
	}

	var config Config

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	return config, nil
}
