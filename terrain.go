package main

import (
	"bytes"
	"fmt"
	"image"
	"net/url"

	"github.com/engelsjk/gmartini"
	"github.com/fogleman/gg"
)

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

	buf := getImage(terrain, mesh, tileSizeInt)

	return buf.Bytes(), nil
}

func getImage(terrain []float32, mesh *gmartini.Mesh, tilesize int) *bytes.Buffer {
	dc := gg.NewContext(tilesize, tilesize)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	drawTriangles(dc, mesh)

	buf := new(bytes.Buffer)
	dc.EncodePNG(buf)
	return buf
}

func drawTriangles(dc *gg.Context, mesh *gmartini.Mesh) {
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
}
