#!/bin/bash

# Function to kill any running process
kill_running_process() {
    # Killing the Go application if it exists
    pid=$(ps aux | grep 'tmp/main' | grep -v 'grep' | awk '{print $2}')
    if [ -n "$pid" ]; then
        echo "Killing existing Go application with PID: $pid"
        kill -9 $pid
    fi
    
    # Killing any running templ generate if exists
    pid=$(ps aux | grep 'templ generate' | grep -v 'grep' | awk '{print $2}')
    if [ -n "$pid" ]; then
        echo "Killing existing templ generate process with PID: $pid"
        kill -9 $pid
    fi
}

# Run 'make css' in the background
make css &

# Stop any existing templ generate process
kill_running_process

# Run 'templ generate' (the main process)
templ generate &

# Build Go application
/usr/bin/go/go/bin/go build -tags=dev -o ./tmp/main .

# Start the Go application in the background with nohup
nohup ./tmp/main

# Ensure Jenkins job doesn't wait for the background processes
echo "App is running independently."
