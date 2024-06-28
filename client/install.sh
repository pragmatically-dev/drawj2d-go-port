#!/bin/bash

BIN_DIR="/home/root/.local/bin/png2rm"
CONFIG_DIR="/home/root/.config/png2rm"
SYSTEMD_DIR="/etc/systemd/system"

CLIENT_SRC="./client"
CONFIG_SRC="./config.yaml"
SERVICE_FILE="png2rm.service"

SERVICE_CONTENT="[Unit]
Description=Client Service Of Png2RmLines
After=home.mount

[Service]
ExecStart=$BIN_DIR/client
Restart=on-failure
RestartSec=5
Enviroment="HOME=/home/root"

[Install]
WantedBy=multi-user.target"

mkdir "$BIN_DIR"
mkdir "$CONFIG_DIR"
mkdir "/home/root/rmdoc"

cp "$CLIENT_SRC" "$BIN_DIR"
if [ $? -ne 0 ]; then
  echo "Error copying client to $BIN_DIR"
  exit 1
fi


cp "$CONFIG_SRC" "$CONFIG_DIR"
if [ $? -ne 0 ]; then
  echo "Error copying config.yaml to $CONFIG_DIR"
  exit 1
fi

echo "Files copied successfully."


read -p "Do you want to edit the config.yaml file? (y/n) " edit_answer

if [[ $edit_answer =~ ^[Yy]$ ]]; then
  nano "$CONFIG_DIR/config.yaml"
  if [ $? -ne 0 ]; then
    echo "Error opening config.yaml with nano"
    exit 1
  fi
fi


echo "$SERVICE_CONTENT" > "$SERVICE_FILE"
mv "$SERVICE_FILE" "$SYSTEMD_DIR"
if [ $? -ne 0 ]; then
  echo "Error creating systemd service file"
  exit 1
fi


systemctl daemon-reload

systemctl enable png2rm.service
 if [ $? -ne 0 ]; then
    echo "Error enabling png2rm service"
    exit 1
  fi

echo "Client service enabled."


systemctl start png2rm.service
  if [ $? -ne 0 ]; then
    echo "Error starting png2rm service"
    exit 1
  fi
echo "png2rm service started."


echo "Installation completed successfully."
