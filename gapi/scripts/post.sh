#!/bin/bash

read -p 'Token: ' token
curl -X POST -d "${2}" -H "Content-Type: application/json" -H "Authorization: Bearer ${token}" ${1}
