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

A [gin](https://github.com/gin-gonic/gin) web API runs in the background on localhost:8080, listening for raster tile requests from the map. The API parses the Z/X/Y coordinate being requested and then forwards that request to the public [Mapbox Terrain-RGB API](https://docs.mapbox.com/help/troubleshooting/access-elevation-data/). Using the Terrain-RGB tile response, an RTIN terrain mesh is created on-the-fly using the [gmartini mesh generator](https://github.com/engelsjk/gmartini/). Mesh triangles are drawn into a PNG file which is then served as the response to the original map raster tile being requested.

