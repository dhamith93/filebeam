#!/bin/bash

input_file="words.txt"
output_file="words.go"
array_name="words"

# Read the word list from the input file
IFS=$'\n' read -d '' -r -a words < "$input_file"

# Generate the Go array
array_definition="package keygen

var $array_name = []string{"

for word in "${words[@]}"; do
  array_definition+="
      \"$word\","
done

array_definition+="
}"

# Write the Go array to the output file
echo "$array_definition" > "$output_file"

echo "Go array '$array_name' generated in '$output_file'"
