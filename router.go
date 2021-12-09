package slippygmartini

import (
	"image"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(slippy *Slippy) {

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}

	slippy.router.Use(cors.New(corsConfig))

	slippy.router.GET("/:format/:z/:x/:y", func(c *gin.Context) {

		z := c.Param("z")
		y := c.Param("y")
		x := c.Param("x")
		format := c.Param("format")

		if format == "vector" {
			slippy.TileSize = 512 // Mapbox GL JS requires vector tile size of 512, but tidwall/mvt uses 256?
		}

		url := slippy.TerrainURL(z, x, y)

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

		mesh, err := slippy.Mesh(img)
		if err != nil {
			c.String(http.StatusInternalServerError, "%s", err.Error())
			return
		}

		switch format {
		case "raster":
			if slippy.Raster(mesh, c); err != nil {
				c.String(http.StatusInternalServerError, "%s", err.Error())
				return
			}
		case "vector":
			if slippy.Vector(mesh, c); err != nil {
				c.String(http.StatusInternalServerError, "%s", err.Error())
				return
			}
		default:
			c.String(http.StatusInternalServerError, "%s", "format request not recognized")
			return
		}
	})
}
