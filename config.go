package slippygmartini

import "os"

type Config struct {
	Port              string
	MapboxAccessToken string
}

func loadConfig() *Config {
	port, _ := os.LookupEnv("PORT")
	mapboxAccessToken, _ := os.LookupEnv("MAPBOX_ACCESS_TOKEN")
	return &Config{
		Port:              port,
		MapboxAccessToken: mapboxAccessToken,
	}
}
