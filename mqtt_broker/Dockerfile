# Use the official Eclipse Mosquitto image
FROM eclipse-mosquitto:latest

# Copy the mosquitto configuration file
COPY mosquitto.conf /mosquitto/config/mosquitto.conf

# Copy the script to generate the password file
COPY generate_password_file.sh /generate_password_file.sh

# Make the script executable
RUN chmod +x /generate_password_file.sh

# Add a debugging step to print contents of /mosquitto/config/
RUN ls -al /mosquitto/config/

# Expose the MQTT port
EXPOSE 1883

# Start Mosquitto and generate password file
CMD sh -c "/generate_password_file.sh && exec mosquitto -c /mosquitto/config/mosquitto.conf"
