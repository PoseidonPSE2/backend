FROM eclipse-mosquitto:openssl

COPY mosquitto.conf /mosquitto/config/mosquitto.conf

CMD ["mosquitto", "-c", "/mosquitto/config/mosquitto.conf"]

EXPOSE 1883
EXPOSE 9001