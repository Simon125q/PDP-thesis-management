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

kill_running_process

# Run 'make css' in the background
make css &

sudo air > output.txt 2>&1

# Run 'templ generate' (the main process)
templ generate &

# Build Go application
/usr/bin/go/go/bin/go build -tags=dev -o ./tmp/main .

# Run the Go app in a detached tmux session
tmux new-session -d -s go_app './tmp/main'

# Log that the app is running in the background inside tmux
echo "App is running in the background inside tmux session 'go_app'."

