#!/bin/bash

# Function to wait for user input before proceeding
pause() {
    read -p "🔹 Press Enter to continue..."
}

# Build the Docker image
echo "🔨 Building Docker image..."
echo "with this command ->docker image build -f Dockerfile -t ascii-art-web-dockerize .
"
docker image build -f Dockerfile -t ascii-art-web-dockerize .
pause

# List all Docker images
echo "📦 Listing Docker images..."
echo "with this command -> docker images
"
docker images
pause

# Run the Docker container in detached mode
echo "🚀 Running the Docker container..."
echo "with this command ->docker container run -p 8080:8080 --detach --name ascii-art ascii-art-web-dockerize
"
docker container run -p 8080:8080 --detach --name ascii-art ascii-art-web-dockerize
pause

# List all containers (running and stopped)
echo "🔍 Listing all containers..."
echo "with this command -> docker ps -a
"
docker ps -a
pause

# Open interactive shell inside the container AND run `ls -l`
echo "🔗 Accessing the container shell..."
echo "with this command -> docker exec -it ascii-art /bin/sh
(sh instead of bash because of Alpine)"
echo "and inside the container we run this command -> ls -l
"
docker exec -it ascii-art /bin/sh -c "ls -l; exec sh"
