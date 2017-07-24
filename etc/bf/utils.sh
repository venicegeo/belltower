#!/bin/bash

assert()
{
	# Function Call:
	# assert "Name of Test" "ARG1" "ARG2" "ARG3" "SKIP?"
	#
	# This function will perform the test [ $ARG1$ARG2$ARG3 ], then
	# display the "Name of Test" in either Green (✓) or Red (✗),
	# depending on the result of the test.
	#
	# If the 5th argument is "skip", then the test will not be
	# performed, and the "Name of Test" will be displayed in Blue (☁)
	#
	# Examples:
	# assert "This tests that integers are equal" 200 -eq 200
	# assert "This tests that integers are not equal" 1 -ne 2
	# assert "This tests that strings are equal" "some_string" == "some_string"
	# assert "This tests that strings not equal" "this_string" != "that_string"
	# assert "This tests that a variable is empty" "" -z "$variable"
	# assert "This tests that a variable is not empty" ! -z "$variable"
	
	RED='\033[0;31m'
	GREEN='\033[0;32m'
	BLUE='\033[0;36m'
	NC='\033[0m'
	
	ASSERTION="$1"
	
	if [ ! -z "$5" ] && [ "skip" == "$5" ]; then
		echo "  $BLUE☁ $ASSERTION$NC"
		(( skips += 1 ))
		return 0
	fi
	
	
	
	LEFT_temp=${2%\"}			# Remove trailing and leading double quotes.
	LEFT=${LEFT_temp#\"}
	if [ -n "$LEFT" ]; then		# Unless empty, add a trailing space.
		LEFT="$LEFT "
	fi
	
	OP=$3
	
	RIGHT_temp=${4%\"}			# Remove trailing and leading double quotes.
	RIGHT=${RIGHT_temp#\"}
	if [ -n "$RIGHT" ]; then	# Unless empty, add a leading space.
		RIGHT=" $RIGHT"
	fi
	
	echo L $LEFT
	echo O $OP
	echo R "${RIGHT}"
	if [ $LEFT$OP$RIGHT ]; then
		echo "  $GREEN✓ $ASSERTION$NC"
		(( passes += 1 ))
	else
		echo "  $RED✗ $ASSERTION$NC"
		(( fails += 1 ))
	fi
}

info()
{
	# This echos the provided text in dark grey.  This uses the
	# same indention as assertions.
	
	NC='\033[0m'
	GREY='\033[0;90m'
	echo "$GREY  ★ $1$NC"
}

http_request()
{
	# Function Call:
	# http_request "METHOD" "https://url.to.hit" "username:password" {"json":"payload"} "HEADER1" "HEADER2" "HEADER3 etc"
	#
	# This function sends the specified request to the provided URL, adding 
	# an auth header if provided one.
	#
	# Running this function will display the http request like so:
	# METHOD: https://url.to.hit
	#   ★ PAYLOAD: {"json":"payload"}
	#   ★ HEADER: "HEADER1"
	#
	# All parameters after the 4th will be interpretted as separate headers;
	# each one will be sent with the -H tag.
	#
	# After this function completes, the variables $body (the response body)
	# and $code (the status code) are available in the calling script.
	
	AUTH="$3"
#	if [ -n "$AUTH" ]; then		# Unless empty, add --user tag.
#		AUTH="--user $AUTH"
#	fi
	
	METHOD="$1"
	URL="$2"
	printf "$METHOD: $URL\n"
	if [ -n "$METHOD" ]; then		# Unless empty, add -X tag.
		METHOD="-X $METHOD"
	fi
	
	HEADERS=""
	for H in "${@:5}"
	do
		HEADERS="$HEADERS -H \"$H\""
		info "HEADER: $H"
	done
	
	PAYLOAD="$4"
	if [ -n "$PAYLOAD" ]; then		# Unless empty, add -d tag.
		info "PAYLOAD: $PAYLOAD"
		PAYLOAD="-d '$PAYLOAD'"
	fi
	
	yy=`curl $METHOD $AUTH $PAYLOAD -s -w ✓%{http_code} $HEADERS $URL`
	echo ----
	echo curl $METHOD $AUTH $PAYLOAD -s -w ✓%{http_code} $HEADERS $URL
	echo ----
	code=${yy##*✓}
	#echo CODE $code
	#echo ----
	body=${yy%✓*}
	#echo BODY $body
	#echo ----

	echo RESPONSE CODE: $code
	echo RESPONSE BODY: $body

	echo vvvvvvvvvvvvv
	echo curl -S -s "$METHOD" $HEADERS -u "$AUTH""" "$PAYLOAD"  $URL
	ret=$(curl -S -s "$METHOD" $HEADERS -u "$AUTH""" "$PAYLOAD"  $URL)
	echo $ret
	echo ^^^^^^^^^^^^^^^^^
}

http_get()
{
	# Function Call:
	# http_get "https://url.to.hit" "username:password"
	
	http_request "GET" "$1" "$2" "" "${@:3}"
}

# This function sends a get request, and assigns $body & $code variables from the response
http_delete()
{	
	# Function Call:
	# http_delete "https://url.to.hit" "username:password"
	
	http_request "DELETE" "$1" "$2" "" "${@:3}"

}

# This function sends a POST request, and assigns $body & $code variables from the response
http_post()
{
	# Function Call:
	# http_get "https://url.to.hit" {"json":"payload"} "username:password"
	#
	# A 'Content-Type: application/json' header is always sent.
	
	http_request "POST" "$1" "$2" "$3" "Content-Type:application/json" "${@:4}"
}

assert_jq_array_contains()
{
	# Function Call:
	# assert "Name of Test" '["jq", "array"]' "string2find" "SKIP?"
	#
	# This function asserts that "string2find" is found in an array
	# created by jq.

	ASSERTION="$1"
	JQ_ARRAY=$2
	TARGET="$3"
	
	found=""
	for ITEM in $JQ_ARRAY
	do
		ITEM=$(strip_jq $ITEM)
		
		if [ -z "$ITEM" ]; then
			continue
		fi

		if [ "$TARGET" == "$ITEM" ]; then
			found="Present"
			break
		fi
	done
	
	assert "$ASSERTION" ! -z "$found" "$4"
}

assert_jq_array_lacks()
{
	# Function Call:
	# assert "Name of Test" '["jq", "array"]' "string2find" "SKIP?"
	#
	# This function asserts that "string2find" is NOT found in an array
	# created by jq.
	
	ASSERTION="$1"
	JQ_ARRAY=$2
	TARGET="$3"
	
	found="Missing"
	for ITEM in $JQ_ARRAY
	do
		ITEM=$(strip_jq $ITEM)
		
		if [ -z "$ITEM" ]; then
			continue
		fi
		
		if [ "$TARGET" == "$ITEM" ]; then
			found=""
			break
		fi
	done
	
	assert "$ASSERTION" ! -z "$found" "$4"
}

strip_jq()
{
	# Function Call:
	# VARIABLE=$(stip_jq $VARIABLE)
	#
	# This function removes extra characters that jq adds to make string comparisons
	# easier, particularly when dealing with elements in an array.
	
	OUTPUT_STRING=${1%]}				# Strip jq stuff (quotes, commas, brackets) from array
	OUTPUT_STRING=${OUTPUT_STRING%,}
	OUTPUT_STRING=${OUTPUT_STRING%\"}
	OUTPUT_STRING=${OUTPUT_STRING#[}
	OUTPUT_STRING=${OUTPUT_STRING#\"}
	
	echo "$OUTPUT_STRING"
}

display_result()
{
	# Display the final result of all tests run thus far.

	RED='\033[0;31m'
	GREEN='\033[0;32m'
	BLUE='\033[0;36m'
	NC='\033[0m'
	
	printf "\nFINAL RESULT\n"
	printf "\n----------------------------------\n\n"
	
	if [ $skips -ne 0 ]; then
		echo "$BLUE ☁☁☁ $skips assertions have been skipped.$NC"
	fi
	
	if [ $fails -eq 0 ]; then
		echo "$GREEN ✓✓✓ Test passed! All $passes assertions have passed.$NC"
		printf "\n----------------------------------\n\n"
		return 0
	else
		echo "$RED ✗✗✗ Test failed! $fails out of $(( $passes + $fails )) assertions failed.$NC"
		printf "\n----------------------------------\n\n"
		return 1
	fi

}
passes=0
fails=0
skips=0
info "Utilities imported!"
