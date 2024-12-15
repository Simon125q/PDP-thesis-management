#!/bin/bash

# Function to kill any running process
kill_running_process() {
    echo "Checking for existing processes..."

    # Killing the Go application if it exists
    pid=$(ps aux | grep 'tmp/main' | grep -v 'grep' | awk '{print $2}')
    if [ -n "$pid" ]; then
        echo "Killing existing Go application with PID: $pid"
        kill -9 $pid
        if [ $? -eq 0 ]; then
            echo "Go application killed successfully."
        else
            echo "Failed to kill Go application with PID: $pid."
        fi
    else
        echo "No Go application running."
    fi

    # Killing any running templ generate if exists
    pid=$(ps aux | grep 'templ generate' | grep -v 'grep' | awk '{print $2}')
    if [ -n "$pid" ]; then
        echo "Killing existing templ generate process with PID: $pid"
        kill -9 $pid
        if [ $? -eq 0 ]; then
            echo "templ generate process killed successfully."
        else
            echo "Failed to kill templ generate process with PID: $pid."
        fi
    else
        echo "No templ generate process running."
    fi
}

kill_running_process

# Run 'make css' in the background
echo "Running 'make css'..."
make css &
make_css_pid=$!
echo "'make css' process started with PID: $make_css_pid"

# Run 'templ generate' (the main process)
echo "Running 'templ generate'..."
templ generate &
templ_generate_pid=$!
echo "'templ generate' process started with PID: $templ_generate_pid"

# Build Go application
echo "Building Go application..."
/usr/bin/go/go/bin/go build -tags=dev -o ./tmp/main .
if [ $? -eq 0 ]; then
    echo "Go application built successfully."
else
    echo "Go application build failed."
    exit 1
fi

# Run the Go app in a detached tmux session
echo "Starting Go application in tmux..."

./tmp/main



# Log that the app is running in the background inside tmux
echo "App is running in the background inside tmux session 'go_app'."
