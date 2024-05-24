#!/bin/bash

# API endpoint URL (replace with your actual URL)
API_ENDPOINT="https://poseidon-backend.fly.dev/refill_stations"

# Predefined RefillStation data (replace with your desired values)
NAME="My Refill Station"
DESCRIPTION="This is a great refill station with clean water."
LATITUDE=49.4478  # Latitude of Kaiserslautern (approximate)
LONGITUDE=7.7753   # Longitude of Kaiserslautern (approximate)
ADDRESS="Fritz-Walter-Platz 1, 67655 Kaiserslautern, Germany"
WATER_SOURCE="Nicequelle"  # Adjust if you have specific information 
OPENING_TIMES="Mon-Fri: 8AM-6PM, Sat: 10AM-4PM"  # Sample opening hours
TYPE="Smart"
OFFERED_WATER_TYPES="Mineral & Tap"  # Adjust if needed
ACTIVE=true

# Construct JSON data string
JSON_DATA="{
  \"name\": \"$NAME\",
  \"description\": \"$DESCRIPTION\",
  \"latitude\": $LATITUDE,
  \"longitude\": $LONGITUDE,
  \"address\": \"$ADDRESS\",
  \"waterSource\": \"$WATER_SOURCE\",
  \"openingTimes\": \"$OPENING_TIMES\",
  \"type\": \"$TYPE\",
  \"offeredWaterTypes\": \"$OFFERED_WATER_TYPES\",
  \"active\": $ACTIVE
}"

# Send POST request using curl
curl -X POST -H "Content-Type: application/json" \
  -d "$JSON_DATA" $API_ENDPOINT

echo ""

# Check for curl exit code (0 indicates success)
if [ $? -eq 0 ]; then
  echo "Refill station creation request sent successfully!"
else
  echo "Error sending request. Please check the API endpoint and curl output."
fi