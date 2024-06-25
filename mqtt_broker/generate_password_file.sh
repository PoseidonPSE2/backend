#!/bin/sh
# Generate the mosquitto password file with retry logic
echo "Creating password file for Mosquitto"
if [ -z "$USERNAME" ] || [ -z "$PASSWORD" ]; then
    echo "USERNAME or PASSWORD environment variables are not set"
    exit 1
fi

# Retry up to 3 times with a delay
retries=3
delay=1 # 1 second delay between retries

for i in $(seq 1 $retries); do
    mosquitto_passwd -b -c /mosquitto/config/mosquitto_passwd $USERNAME $PASSWORD
    if [ $? -eq 0 ]; then
        echo "Password file created successfully."
        chmod 600 /mosquitto/config/mosquitto_passwd
        break
    else
        echo "Attempt $i failed. Retrying in $delay seconds..."
        sleep $delay
    fi
done

# Exit if password file creation fails
if [ $? -ne 0 ]; then
    echo "Failed to create password file after $retries attempts. Exiting."
    exit 1
fi
