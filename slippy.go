package slippygmartini

import (
	"fmt"
	"image"
	"log"
	"net"
	"net/url"

	"github.com/engelsjk/gmartini"
	"github.com/engelsjk/mvt"
	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
)

type Slippy struct {
	config    *Config
	router    *gin.Engine
	MeshError float32
	TileSize  int
}

func New() (*Slippy, error) {

	var tileSize int = 512
	var meshError float32 = 50

	config := loadConfig()
	router := gin.Default()

	s := &Slippy{
		config:    config,
		router:    router,
		TileSize:  tileSize,
		MeshError: meshError,
	}

	SetupRouter(s)

	return s, nil
}

func (s *Slippy) Start() error {
	return s.router.Run(net.JoinHostPort("localhost", s.config.Port))
}

func (s *Slippy) Port() string {
	return s.config.Port
}

func (s *Slippy) TerrainURL(z, x, y string) string {

	var dpi string
	if s.TileSize == 512 {
		dpi = "@2x"
	}

	u := url.URL{
		Scheme: "https",
		Host:   "api.mapbox.com",
		Path:   fmt.Sprintf("v4/mapbox.terrain-rgb/%s/%s/%s%s.pngraw", z, x, y, dpi),
	}

	q := u.Query()
	q.Set("access_token", s.config.MapboxAccessToken)
	u.RawQuery = q.Encode()

	return u.String()
}

func (s *Slippy) Mesh(img image.Image) (*gmartini.Mesh, error) {

	terrain, err := gmartini.DecodeElevation(img, "mapbox", true)
	if err != nil {
		return nil, err
	}

	martini, err := gmartini.New(gmartini.OptionGridSize(int32(s.TileSize) + 1))
	if err != nil {
		return nil, err
	}

	tile, err := martini.CreateTile(terrain)
	if err != nil {
		return nil, err
	}

	return tile.GetMesh(gmartini.OptionMaxError(s.MeshError)), nil
}

func (s *Slippy) Raster(mesh *gmartini.Mesh, c *gin.Context) error {

	dc := gg.NewContext(s.TileSize, s.TileSize)

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

	c.Writer.Header().Set("Content-Type", "image/png")
	return dc.EncodePNG(c.Writer)
}

func (s *Slippy) Vector(mesh *gmartini.Mesh, c *gin.Context) error {

	log.Println("vector")

	var tile mvt.Tile
	l := tile.AddLayer("gmartini")
	for i := 0; i < (len(mesh.Triangles) - 3); i += 3 {
		a, b, c := mesh.Triangles[i], mesh.Triangles[i+1], mesh.Triangles[i+2]
		ax, ay := float64(mesh.Vertices[2*a]), float64(mesh.Vertices[2*a+1])
		bx, by := float64(mesh.Vertices[2*b]), float64(mesh.Vertices[2*b+1])
		cx, cy := float64(mesh.Vertices[2*c]), float64(mesh.Vertices[2*c+1])
		f := l.AddFeature(mvt.Polygon)
		f.MoveTo(ax, ay)
		f.LineTo(bx, by)
		f.LineTo(cx, cy)
		f.LineTo(ax, ay)
		f.ClosePath()
	}

	c.Data(200, "application/vnd.mapbox-vector-tile", tile.Render())
	return nil
}
