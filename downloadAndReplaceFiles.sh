#!/bin/bash

mkdir -p samples

image_urls=(
  #add data for the repos to be files to be downloaded
)

for url in "${image_urls[@]}"; do
  file_name=$(basename "$url")
  curl -sSfL -o "samples/$file_name" "$url"
done
