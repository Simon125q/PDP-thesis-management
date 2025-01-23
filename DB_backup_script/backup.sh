#!/bin/bash

# Log file to store the output of the script
LOG_FILE="/var/log/run_python_program.log"

# Path to the Python interpreter
PYTHON="/usr/bin/python3"

# Path to the Python program
SCRIPT="/path/to/your/backupScript.py"

echo "$(date): Starting the Python program..." >> "$LOG_FILE"

# Run the Python program and log output and errors
$PYTHON "$SCRIPT" >> "$LOG_FILE" 2>&1

echo "$(date): Python program execution completed." >> "$LOG_FILE"

