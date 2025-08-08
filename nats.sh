docker run -d \
  --name nats \
  --restart unless-stopped \
  -p 4222:4222 \
  -p 8222:8222 \
  -p 6222:6222 \
  -v "${PWD}/nats-server.conf:/etc/nats/nats-server.conf" \
  -v nats-data:/data \
  nats:latest \
  -c /etc/nats/nats-server.conf