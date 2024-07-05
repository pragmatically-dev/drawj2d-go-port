#!/bin/bash

BIN_DIR="/home/root/.local/bin/drawj2d-go"
SYSTEMD_DIR="/etc/systemd/system"

CLIENT_SRC="./drawj2d-go"

SERVICE_FILE="drawj2d-go.service"

SERVICE_CONTENT="[Unit]
Description=Png to rmlines conversion
After=home.mount

[Service]
ExecStart=$BIN_DIR/drawj2d-go
Restart=on-failure
RestartSec=5
Enviroment="HOME=/home/root"

[Install]
WantedBy=multi-user.target"

mkdir "$BIN_DIR"
mkdir "$CONFIG_DIR"


cp "$CLIENT_SRC" "$BIN_DIR"
if [ $? -ne 0 ]; then
  echo "Error copying client to $BIN_DIR"
  exit 1
fi


echo "Files copied successfully."


echo "$SERVICE_CONTENT" > "$SERVICE_FILE"
mv "$SERVICE_FILE" "$SYSTEMD_DIR"
if [ $? -ne 0 ]; then
  echo "Error creating systemd service file"
  exit 1
fi


systemctl daemon-reload

systemctl enable drawj2d-go.service
 if [ $? -ne 0 ]; then
    echo "Error enabling drawj2d-go service"
    exit 1
  fi

echo "Client service enabled."


systemctl start drawj2d-go.service
  if [ $? -ne 0 ]; then
    echo "Error starting drawj2d-go service"
    exit 1
  fi
echo "drawj2d-go service started."


echo "Installation completed successfully."