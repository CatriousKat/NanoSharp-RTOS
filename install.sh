#!/bin/bash
set -e

echo "=== NanoSharp Universal Microcontroller Flasher ==="

if [ ! -f "interpreter.go" ]; then
    echo "[ERROR] interpreter.go not found in the current directory!"
    exit 1
fi

echo -e "\nSelect your target microcontroller:"
echo "  1) esp32 (Generic ESP32)"
echo "  2) esp32c3 (ESP32-C3)"
echo "  3) pico (Raspberry Pi Pico)"
echo "  4) pico2 (Raspberry Pi Pico 2)"
echo "  5) arduino-nano33 (Arduino Nano 33 IoT)"
echo "  6) Custom target name"

read -p "Enter your choice (1-6): " choice

case $choice in
    1) target="esp32" ;;
    2) target="esp32c3" ;;
    3) target="pico" ;;
    4) target="pico2" ;;
    5) target="arduino-nano33" ;;
    6) read -p "Enter custom TinyGo target name: " target ;;
    *) target="esp32" ;;
esac

echo -e "\n[INFO] Compiling and flashing interpreter.go for target: $target..."
tinygo flash -target="$target" interpreter.go

echo "[SUCCESS] NanoSharp successfully flashed to $target!"
echo "=== Flash Complete ==="