#!/bin/bash

APP_NAME="DeskNotifier"
OUTPUT_DIR="output"
TARGET_DIR="/opt/$APP_NAME"
SERVICE_FILE="/etc/systemd/system/$APP_NAME.service"

# Create output directory if it doesn't exist
if [ ! -d "$OUTPUT_DIR" ]; then
    mkdir "$OUTPUT_DIR"
fi

# Set environment variables for cross-compilation
export GOOS=linux
export GOARCH=arm
export GOARM=7

# Build the application for Raspbian
echo "Building app for Raspbian..."
go build -o "$OUTPUT_DIR/$APP_NAME" ./cmd/main.go

if [ $? -ne 0 ]; then
    echo "Error while building app!"
    exit 1
fi

# Create install.sh script
echo "Creating install.sh..."

cat <<EOL > "$OUTPUT_DIR/install.sh"
#!/bin/bash

APP_NAME="$APP_NAME"
TARGET_DIR="$TARGET_DIR"
SERVICE_FILE="$SERVICE_FILE"

echo "Using current directory: \$TARGET_DIR"

if [ ! -f "\$TARGET_DIR/\$APP_NAME" ]; then
    echo "Error: \$APP_NAME not found in current directory!"
    exit 1
fi

echo "Setting execute permissions on \$TARGET_DIR/\$APP_NAME..."
chmod +x "\$TARGET_DIR/\$APP_NAME"

echo "Creating service file..."
sudo bash -c "cat > \$SERVICE_FILE << EOF
[Unit]
Description=Standing Desk Notifier by Maciej Rosiak
After=network.target

[Service]
ExecStart=\$TARGET_DIR/\$APP_NAME
Restart=always
User=pi
Group=pi
WorkingDirectory=\$TARGET_DIR
Environment=GO_ENV=production

[Install]
WantedBy=multi-user.target
EOF"

echo "Setting permissions for \$SERVICE_FILE..."
sudo chmod 644 \$SERVICE_FILE

echo "Reloading systemd daemon..."
sudo systemctl daemon-reload

echo "Enabling service..."
sudo systemctl enable \$APP_NAME

echo "Starting service..."
sudo systemctl start \$APP_NAME

echo "Checking service status..."
sudo systemctl status \$APP_NAME

echo "✔️ Installation completed!"
EOL

# Give execute permission to install.sh
chmod +x "$OUTPUT_DIR/install.sh"

echo "install.sh script has been created and made executable in the $OUTPUT_DIR folder."
