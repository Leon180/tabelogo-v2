#!/bin/bash

echo "=== Docker Complete Removal Script ==="
echo "This script will completely remove Docker from your Mac"
echo ""

# Stop Docker Desktop
echo "Step 1: Quitting Docker Desktop..."
osascript -e 'quit app "Docker"' 2>/dev/null
sleep 3

# Kill Docker processes
echo "Step 2: Killing Docker processes..."
pkill -9 -f "Docker" 2>/dev/null
sleep 2

# Remove Docker.app
echo "Step 3: Removing Docker.app..."
sudo rm -rf /Applications/Docker.app

# Remove Docker CLI tools
echo "Step 4: Removing Docker CLI tools..."
sudo rm -f /usr/local/bin/docker
sudo rm -f /usr/local/bin/docker-compose
sudo rm -f /usr/local/bin/docker-credential-desktop
sudo rm -f /usr/local/bin/docker-credential-ecr-login
sudo rm -f /usr/local/bin/docker-credential-osxkeychain
sudo rm -f /usr/local/bin/kubectl
sudo rm -f /usr/local/bin/kubectl.docker
sudo rm -f /usr/local/bin/com.docker.cli

# Remove Docker helper tools
echo "Step 5: Removing Docker helper tools..."
sudo rm -f /Library/PrivilegedHelperTools/com.docker.vmnetd
sudo rm -f /Library/LaunchDaemons/com.docker.vmnetd.plist

# Remove Docker data directories
echo "Step 6: Removing Docker data directories..."
rm -rf ~/Library/Containers/com.docker.docker
rm -rf ~/Library/Application\ Support/Docker\ Desktop
rm -rf ~/Library/Group\ Containers/group.com.docker
rm -rf ~/Library/Preferences/com.docker.docker.plist
rm -rf ~/Library/Preferences/com.electron.docker-frontend.plist
rm -rf ~/Library/Saved\ Application\ State/com.electron.docker-frontend.savedState
rm -rf ~/Library/Logs/Docker\ Desktop
rm -rf ~/.docker

# Remove Docker caches
echo "Step 7: Removing Docker caches..."
rm -rf ~/Library/Caches/com.docker.docker
rm -rf ~/Library/Caches/Docker\ Desktop

echo ""
echo "✅ Docker has been completely removed!"
echo ""
echo "To verify, run:"
echo "  which docker"
echo "  ls /Applications/Docker.app"
echo ""
echo "Both should return 'not found' or similar."
