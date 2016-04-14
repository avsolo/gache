#!/bin/bash

# Simple shortcut to check telnet response
# Format of call:
#   ./telnet.sh <ip> <port> <command>

exec 3<>/dev/tcp/$1/$2
echo -en "${@:3}\n" >&3
RESPONSE="`cat <&3`"
echo "Server response is: $RESPONSE"
