#!/bin/bash -e

# install jq

. ./utils.sh	# Import functions like assert and http_get
. ~/.bf-vars.sh		# Import variables needed for these tests

domain=https://bf-api.geointservices.io

echo -------------------------------------------------------------
ret=$(curl -S -s -X GET -u $auth:""  ${domain}/v0/algorithm)
service_id=$(echo $ret | jq -r '.algorithms[0].service_id')
echo SERVICE ID: $service_id
echo -------------------------------------------------------------


####echo -------------------------------------------------------------
#####http_get "${http}bf-api.$domain/v0/user" "$auth"
#assert "Should receive a 200" 200 -eq "$code"
####catalog_url="$(echo $body | jq -r '.services.catalog')"
#assert "The catalog url should exist" "null" != "$catalog_url"
###echo CATALOG URL: $catalog_url
###echo -------------------------------------------------------------



echo -------------------------------------------------------------
ret=$(curl -S -s -X GET -u $auth:""  ${domain}/v0/algorithm)
echo $ret
echo -------------------------------------------------------------


echo -------------------------------------------------------------
payload='{
	"algorithm_id": "'$service_id'",
	"scene_id": "landsat:'$selected_image'",
	"name": "mpg97",
	"planet_api_key": "'$planet_key'"
}'
#echo curl -S -s -X POST -u $auth:"" -d "$payload" $domain/v0/job
ret=$(curl -S -s -X POST -u $auth:"" -d "$payload" -H "Content-Type: application/json" $domain/v0/job)
job_id="$(echo $ret | jq -r '.job.id')"
echo JOB ID: $job_id
echo -------------------------------------------------------------

for i in {1..20}
do
	echo -------------------------------------------------------------
	ret=$(curl -S -s -X GET -u $auth:""  ${domain}/v0/job/${job_id})
	#http_get "${http}bf-api.$domain/v0/job/$job_id" "$auth"
	status="$(echo $ret | jq -r '.job.properties.status')"
	if [ "Success" == "$status" ]; then
		i=0 # Pass, in case on last iteration.
		break
	else
		info "Job status is $status, trying again.  Attempt $i of 20."
		sleep 30
	fi
	echo -------------------------------------------------------------
done
assert "Job should be successful" 20 -ne "$i"


#http_get "${http}bf-api.$domain/v0/job" "$auth"
#assert "Should receive a 200" 200 -eq "$code"
assert_jq_array_contains "Job should be in list" "$(echo $body | jq -r '[.jobs.features[] | .id]')" "$job_id"

body=$(curl -S -s -X GET -u $auth:""  ${domain}/v0/job/by_scene/landsat:$selected_image)
#http_get "${http}bf-api.$domain/v0/job/by_scene/landsat:$selected_image" "$auth"
#assert "Should receive a 200" 200 -eq "$code"
assert_jq_array_contains "Job should be in jobs listed by scene" "$(echo $body | jq -r '[.jobs.features[] | .id]')" "$job_id"
for IMAGE in $(echo $body | jq -r '[.jobs.features[] | .properties.scene_id]')
do
	IMAGE=$(strip_jq $IMAGE)
	
	if [ -z "$IMAGE" ]; then
		continue
	fi
	
	assert "Image in list has correct scene" "landsat:$selected_image" == "$IMAGE"
done

#http_delete "${http}bf-api.$domain/v0/job/$job_id" "$auth"
