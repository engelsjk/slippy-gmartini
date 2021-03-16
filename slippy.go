package slippygmartini

import (
	"bytes"
	"fmt"
	"image"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/engelsjk/gmartini"
	"github.com/fogleman/gg"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

func Run() {
	config := loadConfig()

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}

	r.Use(cors.New(corsConfig))

	r.GET("/raster/:z/:x/:y", func(c *gin.Context) {

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

type Params struct {
	X         string
	Y         string
	Z         string
	MeshError float64
	TileSize  int
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

func getTerrainURL(token string, z, x, y string, tilesize int) string {
	var dpi string
	if tilesize == 512 {
		dpi = "@2x"
	}

	u := url.URL{
		Scheme: "https",
		Host:   "api.mapbox.com",
		Path:   fmt.Sprintf("v4/mapbox.terrain-rgb/%s/%s/%s%s.pngraw", z, x, y, dpi),
	}

	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()

	return u.String()
}

func getTerrainTile(img image.Image, meshErrInt float64, tileSizeInt int) ([]byte, error) {
	terrain, err := gmartini.DecodeElevation(img, "mapbox", true)
	if err != nil {
		return nil, err
	}

	martini, err := gmartini.New(gmartini.OptionGridSize(int32(tileSizeInt) + 1))
	if err != nil {
		return nil, err
	}

	tile, err := martini.CreateTile(terrain)
	if err != nil {
		return nil, err
	}

	mesh := tile.GetMesh(gmartini.OptionMaxError(float32(meshErrInt)))

	return getRaster(terrain, mesh, tileSizeInt)
}

func getRaster(terrain []float32, mesh *gmartini.Mesh, tilesize int) ([]byte, error) {
	dc := gg.NewContext(tilesize, tilesize)

	dc.ClearPath()
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(0.1)
	for i := 0; i < (len(mesh.Triangles) - 3); i += 3 {
		a, b, c := mesh.Triangles[i], mesh.Triangles[i+1], mesh.Triangles[i+2]
		ax, ay := float64(mesh.Vertices[2*a]), float64(mesh.Vertices[2*a+1])
		bx, by := float64(mesh.Vertices[2*b]), float64(mesh.Vertices[2*b+1])
		cx, cy := float64(mesh.Vertices[2*c]), float64(mesh.Vertices[2*c+1])
		dc.MoveTo(ax, ay)
		dc.LineTo(bx, by)
		dc.LineTo(cx, cy)
		dc.LineTo(ax, ay)
	}
	dc.Stroke()

	buf := new(bytes.Buffer)
	err := dc.EncodePNG(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getVector(terrain []float32, mesh *gmartini.Mesh, tilesize int) []byte {
	return nil
}
