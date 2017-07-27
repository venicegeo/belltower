#!/bin/bash

# usage:
#  ./beachfront.sh [image-id]
#

# sets $domain, $auth, $planet_key, $selected_image
. ~/.bf-vars.sh

url=https://bf-api.$domain



for i in {1..20}
do
	echo -------------------------------------------------------------
	body=$(curl -S -s -X GET -u $auth:""  ${url}/v0/job)
	#echo $body > x

	status="$(echo $body | jq -r '.jobs.features[].properties.status + "\t" +  .jobs.features[].properties.scene_id')"
	#if [ "Success" == "$status" ]; then
	#	echo STATUS: $status
	#	i=0 # Pass, in case on last iteration.
	#	break
	#else
		echo "$status" | grep -v Success |  grep -v Error | grep -v "Timed Out"
		echo "$status" | wc -l
		sleep 3
	#fi
	echo -------------------------------------------------------------
	echo
done

echo -------------------------------------------------------------
curl -S -s -X GET -u $auth:""  ${url}/v0/job/${job_id} > $job_id.json
curl -S -s -X GET -u $auth:""  ${url}/v0/job/${job_id}.geojson > $job_id.geojson
ls -l $job_id.*
echo -------------------------------------------------------------
echo
