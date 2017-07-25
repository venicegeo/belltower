#!/bin/bash -e

. ~/.bf-vars.sh		# Import variables needed for these tests

broker_url=https://bf-ia-broker.$domain

body=$(curl -S -s -X GET -u $auth:""  ${broker_url}/planet/landsat/${selected_image}?PL_API_KEY=${planet_key})
red=$(echo $body | jq -r '.properties.bands.red')
green=$(echo $body | jq -r '.properties.bands.green')
blue=$(echo $body | jq -r '.properties.bands.blue')
echo RED: $red...
curl -S -s -X GET $red > red.tif
echo GREEN: $green...
curl -S -s -X GET $green > green.tif
echo BLUE:  $blue
curl -S -s -X GET $blue > blue.tif
