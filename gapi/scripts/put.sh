#!/bin/bash

read -p 'Token: ' token

curl -X PUT -d "${2}" -H "Content-Type: application/json" -H "Authorization: Bearer ${token}" "${1}"
