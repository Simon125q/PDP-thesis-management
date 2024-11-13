#!/bin/bash

# to run the app on linux run ./runLinux.sh

# Run 'make css' in the background
nohup make css &

# Run 'air' in the background
nohup air &

# Run 'templ generate --watch --proxy=http://localhost:3000' in the background
nohup templ generate --watch --proxy=http://localhost:3000 &


