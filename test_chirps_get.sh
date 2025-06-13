#!/bin/bash

# Reset the database
curl -s -X POST http://localhost:8080/admin/reset

# Get user's id
user_response=$(curl -s -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "saul@bettercall.com"}')


# Extract user's id from response
echo
userID1=$(echo "$user_response" | jq -r '.id')
echo "Created user with ID: $userID1"

# Create a chirp with the user's ID
chirp_response=$(curl -s -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -d "{\"body\": \"If you're committed enough, you can make any story work.\", \"user_id\": \"$userID1\"}")

echo "Chirp response:"
echo "$chirp_response"

# Create another chirp with the same user's id
chirp_response_2=$(curl -s -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -d "{\"body\": \"I once told a woman I was Kevin Costner, and it worked because I believed it.\", \"user_id\": \"$userID1\"}")

echo "Chirp response 2:"
echo "$chirp_response_2"

# Get all chirps
chirps=$(curl -s -X GET http://localhost:8080/api/chirps)
echo "Chirps:"
echo "$chirps"