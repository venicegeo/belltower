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
	
	if [ $LEFT$OP$RIGHT ]; then
		echo "  $GREEN✓ $ASSERTION$NC"
		(( passes += 1 ))
	else
		echo "  $RED✗ $ASSERTION$NC"
		(( fails += 1 ))
	fi
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
