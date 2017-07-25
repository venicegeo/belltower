#!/bin/bash -e

# brew install gdal-20
# export PATH="/usr/local/opt/gdal2/bin:$PATH"
# brew install gdal2-python
# export PATH="/usr/local/opt/gdal2-python/bin:$PATH"

rm -f rgb.tif
gdal_merge.py -separate -o rgb.tif red.tif green.tif blue.tif
gdal_rasterize -b 1 -b 2 -b 3 -burn 65535 -burn 65535 -burn 65535 in.geojson rgb.tif
