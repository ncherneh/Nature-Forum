#!/bin/bash

# Set image name
forum_docker="my_app"

# Stop the running container
docker container stop $(docker ps -q -f ancestor=$forum_docker)

# Remove the stopped container
docker container rm $(docker ps -a -q -f ancestor=$forum_docker)

# Remove the Docker image
docker image rm -f $forum_docker

# docker system prune -f