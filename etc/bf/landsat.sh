#!/bin/bash -e

# install jq

. ./utils.sh	# Import functions like assert and http_get
. ~/.bf-vars.sh		# Import variables needed for these tests

domain=https://bf-api.int.geointservices.io
broker_domain=https://bf-ia-broker.int.geointservices.io

body=$(curl -S -s -X GET -u $auth:""  ${broker_domain}/planet/landsat/${selected_image}?PL_API_KEY=${planet_key})
red=$(echo $body | jq -r '.properties.bands.red')
green=$(echo $body | jq -r '.properties.bands.green')
blue=$(echo $body | jq -r '.properties.bands.blue')
echo RED: $red...
curl -S -s -X GET $red > red.tif
echo GREEN: $green...
curl -S -s -X GET $green > green.tif
echo BLUE:  $blue
curl -S -s -X GET $blue > blue.tif

# brew install gdal-20
# export PATH="/usr/local/opt/gdal2/bin:$PATH"
# brew install gdal2-python
# export PATH="/usr/local/opt/gdal2-python/bin:$PATH"
#
# rm rgb.tif
# gdal_merge.py -separate -o rgb.tif red.tif green.tif blue.tif
# gdal_rasterize -b 1 -b 2 -b 3 -burn 65535 -burn 65535 -burn 65535 in.geojson rgb.tif