#!/bin/bash

APP_NAME="DeskNotifier"

MAIN_FILE="./cmd/main.go"

TARGET_DIR="/opt/${APP_NAME}"

SERVICE_FILE="/etc/systemd/system/${APP_NAME}.service"

echo "Building app.."
go build -o ${APP_NAME} ${MAIN_FILE}

if [ $? -ne 0 ]; then
    echo "Error while building app!"
    exit 1
fi

echo "Creating ${TARGET_DIR} directory..."
sudo mkdir -p ${TARGET_DIR}

echo "Moving app to ${TARGET_DIR}..."
sudo mv ${APP_NAME} ${TARGET_DIR}

if [ $? -ne 0 ]; then
    echo "Error while moving app!"
    exit 1
fi

echo "Creating service file..."
sudo bash -c "cat > ${SERVICE_FILE}" << EOL
[Unit]
Description=Standing Desk Notifier by Maciej Rosiak
After=network.target

[Service]
ExecStart=${TARGET_DIR}/${APP_NAME}
Restart=always
User=pi
Group=pi
WorkingDirectory=${TARGET_DIR}
Environment=GO_ENV=production

[Install]
WantedBy=multi-user.target
EOL

echo "Setting permissions for ${SERVICE_FILE}..."
sudo chmod 644 ${SERVICE_FILE}

echo "Reloading systemd daemon..."
sudo systemctl daemon-reload

echo "Enabling service..."
sudo systemctl enable ${APP_NAME}

echo "Starting service..."
sudo systemctl start ${APP_NAME}

echo ""
echo "✔️ Installation completed!"
