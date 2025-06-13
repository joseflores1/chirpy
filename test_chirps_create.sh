#!/bin/bash

# Reset database
curl -s -X POST http://localhost:8080/admin/reset

# Create the user and capture the response
user_response=$(curl -s -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "saul@bettercall.com"}')

# Extract user ID from the response
echo 
userID1=$(echo "$user_response" | jq -r '.id')
echo "Created user with ID: $userID1"

# Create a chirp with the user's ID
chirp_response=$(curl -s -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -d "{\"body\": \"If you're committed enough, you can make any story work.\", \"user_id\": \"$userID1\"}")
echo "Chirp response:"
echo "$chirp_response"