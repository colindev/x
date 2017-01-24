#!/bin/bash

read -p 'Token: ' token
curl -H "Content-Type: application/json" -H "Authorization: Bearer ${token}" "${1}"
