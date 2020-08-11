# slippy-rtin

![](/images/atlantic.png)

A visual experiment using RTIN terrain meshes in a slippy tile map.

## The How

A Mapbox GL JS map runs locally (```localhost:5500```) with a raster layer sourced by a locally-hosted slippy tile service running on another port.

```javascript
'sources': {
    'raster-tiles': {
      'type': 'raster',
      'tiles': [
        'http://localhost:8080/{z}/{x}/{y}'
      ],
      'tileSize': 512
    }
  }
```

A [gin](https://github.com/gin-gonic/gin) web API runs in the background on localhost:8080, listening for raster tile requests from the map. The API parses the Z/X/Y tile coordinate being requested and then forwards that request to the public [Mapbox Terrain-RGB API](https://docs.mapbox.com/help/troubleshooting/access-elevation-data/). 

The Terrain-RGB tile that is returned from the Mapbox API is then used to create an RTIN terrain mesh using the [gmartini mesh generator](https://github.com/engelsjk/gmartini/). Finally, mesh triangles from the mesh are drawn into a PNG file which is served as a response to the original map raster tile that was requested.

Think of it as an on-the-fly slippy tile transmogrifying service.

