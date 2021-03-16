# slippy-rtin

![](/images/atlantic.png)

A visual experiment using RTIN terrain meshes in a slippy tile map.

## The How

Think of this as an on-the-fly slippy tile transmogrifying experiment.

A slippy tile API runs locally on one port, listening for raster tile requests from a web map. The API parses the Z/X/Y tile coordinate being requested and forwards that tile request to the [Mapbox Terrain-RGB API](https://docs.mapbox.com/help/troubleshooting/access-elevation-data/). The Terrain-RGB tile that is returned from the Mapbox API is then used to create an RTIN terrain mesh using the [gmartini mesh generator](https://github.com/engelsjk/gmartini/). Finally, flattened mesh triangles are drawn into a transparent PNG file which is then served in response to the slippy tile API request.

A web map runs locally on another port, serving a basemap layer and a raster layer of gmartini slippy tiles.

```javascript
map.addSource('gmartini', {
  'type': 'raster',
  'tiles': [
    'http://localhost:8080/raster/{z}/{x}/{y}'
  ],
  'tileSize': 512
});

map.addLayer({
  'id': 'rtin',
  'type': 'raster',
  'source': 'gmartini',
  'minzoom': 0,
  'maxzoom': 22
});
```

## Run

First, start up the slippy tile service. You must have both a desired port and a Mapbox token (to access the Terrain-RGB api) specified in a .env file.

```bash
go run cmd/main.go
```

Next, use a simple static web server, like [m3ng9i/ran](https://github.com/m3ng9i/ran), to serve the web map. Edit ```web/index.html``` to include a Mapbox token and to make sure that the tile source url contains the same port that the slippy tile service is running on.

```bash
ran -r web -p 5500
```

Now, just stare at some triangles...

![](/images/slippy.png)
