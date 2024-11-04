#!/bin/bash

# to run the app on linux run ./runLinux.sh

# Run 'make css' in the background
make css &

# Run 'air' in the background
air &

# Run 'templ generate --watch --proxy=http://localhost:3000' in the background
templ generate --watch --proxy=http://localhost:3000 &


