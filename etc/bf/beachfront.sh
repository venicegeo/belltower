#!/bin/bash -e

# install jq

. ./utils.sh	# Import functions like assert and http_get
. ~/.bf-vars.sh		# Import variables needed for these tests

domain=https://bf-api.int.geointservices.io
broker_domain=https://bf-ia-broker.int.geointservices.io


echo -------------------------------------------------------------
ret=$(curl -S -s -X GET -u $auth:""  ${domain}/v0/algorithm)
service_id=$(echo $ret | jq -r '.algorithms[0].service_id')
echo SERVICE ID: $service_id
echo -------------------------------------------------------------
echo

echo -------------------------------------------------------------
payload='{
	"algorithm_id": "'$service_id'",
	"scene_id": "landsat:'$selected_image'",
	"name": "mpg97",
	"planet_api_key": "'$planet_key'"
}'
ret=$(curl -S -s -X POST -u $auth:"" -d "$payload" -H "Content-Type: application/json" $domain/v0/job)
job_id="$(echo $ret | jq -r '.job.id')"
echo JOB ID: $job_id
echo -------------------------------------------------------------
echo

for i in {1..20}
do
	echo -------------------------------------------------------------
	body=$(curl -S -s -X GET -u $auth:""  ${domain}/v0/job/${job_id})
	#echo $body
	status="$(echo $body | jq -r '.job.properties.status')"
	if [ "Success" == "$status" ]; then
		echo STATUS: $status
		i=0 # Pass, in case on last iteration.
		break
	else
		echo "$i: $status"
		sleep 30
	fi
	echo -------------------------------------------------------------
	echo
done

echo -------------------------------------------------------------
curl -S -s -X GET -u $auth:""  ${domain}/v0/job/${job_id} > $job_id.json
curl -S -s -X GET -u $auth:""  ${domain}/v0/job/${job_id}.geojson > $job_id.geojson
ls -l $job_id.*
echo -------------------------------------------------------------
echo
