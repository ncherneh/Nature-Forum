#!/bin/bash

# Set image name
forum_docker="my_app"

# Build an image from Dockerfile
docker build -t $forum_docker .

# Start the container from the image
docker run -p 8083:8083 $forum_docker