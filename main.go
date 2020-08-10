package main

import (
	"fmt"
	"image"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	MapboxAccessToken string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

type Params struct {
	X         string
	Y         string
	Z         string
	MeshError float64
	TileSize  int
}

func main() {
	config := loadConfig()

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://127.0.0.1:5500"}

	r.Use(cors.New(corsConfig))

	r.GET("/:z/:x/:y", func(c *gin.Context) {

		params, err := parseParams(c)
		if err != nil {
			c.String(http.StatusBadRequest, "%s", err.Error())
			return
		}

		url := getTerrainURL(config.MapboxAccessToken, params.Z, params.X, params.Y, params.TileSize)
		resp, err := http.Get(url)
		if err != nil {
			c.String(http.StatusInternalServerError, "%s", err.Error())
			return
		}
		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "%s", err.Error())
			return
		}

		bytes, err := getTerrainTile(img, params.MeshError, params.TileSize)
		if err != nil {
			c.String(http.StatusInternalServerError, "%s", err.Error())
			return
		}

		c.Data(http.StatusOK, "png", bytes)
	})

	r.Run(net.JoinHostPort("localhost", config.Port))
}

func loadConfig() *Config {
	port, _ := os.LookupEnv("PORT")
	mapboxAccessToken, _ := os.LookupEnv("MAPBOX_ACCESS_TOKEN")
	return &Config{
		Port:              port,
		MapboxAccessToken: mapboxAccessToken,
	}
}

func parseParams(c *gin.Context) (*Params, error) {

	z := c.Param("z")
	y := c.Param("y")
	x := c.Param("x")

	tilesizeStr := c.DefaultQuery("tile", "512")
	if !(tilesizeStr == "256" || tilesizeStr == "512") {
		return nil, fmt.Errorf("tile must be either 256 or 512")
	}
	tileSize, err := strconv.Atoi(tilesizeStr)
	if err != nil {
		return nil, fmt.Errorf("tile must be an integer")
	}

	meshErrStr := c.DefaultQuery("mesh", "25")
	meshErr, err := strconv.ParseFloat(meshErrStr, 32)
	if err != nil {
		return nil, fmt.Errorf("mesh error must be a number")
	}

	return &Params{
		X:         x,
		Y:         y,
		Z:         z,
		TileSize:  tileSize,
		MeshError: meshErr,
	}, nil
}
