# NanoSharp Universal Microcontroller Flasher (PowerShell)
$ErrorActionPreference = "Stop"

Write-Host "=== NanoSharp Universal Microcontroller Flasher ===" -ForegroundColor Cyan

if (-not (Test-Path "interpreter.go")) {
    Write-Host "[ERROR] interpreter.go not found in the current directory!" -ForegroundColor Red
    exit 1
}

if (-not (Test-Path "script.ns")) {
    Write-Host "[ERROR] script.ns not found! The script needs to be present so it can be embedded into the firmware." -ForegroundColor Red
    exit 1
}

# Automatically locate true TinyGo root directory
if (-not $env:TINYGOROOT) {
    $wingetTinyGo = Get-ChildItem "$env:LOCALAPPDATA\Microsoft\WinGet\Packages" -Filter "*tinygo*" -Directory -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($wingetTinyGo) {
        $foundRoot = Get-ChildItem $wingetTinyGo.FullName -Filter "tinygo" -Directory -Recurse -ErrorAction SilentlyContinue | Select-Object -First 1
        if ($foundRoot -and (Test-Path (Join-Path $foundRoot.FullName "src"))) {
            $env:TINYGOROOT = $foundRoot.FullName
        }
    }
}

if (-not $env:TINYGOROOT) {
    $defaultPaths = @(
        "C:\Program Files\tinygo",
        "$env:LOCALAPPDATA\Programs\tinygo",
        "C:\tinygo"
    )
    foreach ($path in $defaultPaths) {
        if (Test-Path (Join-Path $path "src")) {
            $env:TINYGOROOT = $path
            break
        }
    }
}

if (-not $env:TINYGOROOT) {
    Write-Host "[WARNING] Could not automatically locate TinyGo root directory." -ForegroundColor Yellow
    $env:TINYGOROOT = Read-Host "Please enter the full path to your main TinyGo folder (the folder that contains the 'src' subfolder)"
}

[Environment]::SetEnvironmentVariable("TINYGOROOT", $env:TINYGOROOT, "Process")
Write-Host "[INFO] Using TINYGOROOT: $env:TINYGOROOT" -ForegroundColor DarkGray

Write-Host "`nSelect your target development board:" -ForegroundColor Yellow
Write-Host "  1) esp32-mini32 (Generic ESP32 / Mini32 Board)"
Write-Host "  2) xiao-esp32c3 (Seeed Studio XIAO ESP32-C3)"
Write-Host "  3) pico (Raspberry Pi Pico)"
Write-Host "  4) pico2 (Raspberry Pi Pico 2)"
Write-Host "  5) arduino-nano33 (Arduino Nano 33 IoT)"
Write-Host "  6) Custom target name"

$choice = Read-Host "`nEnter your choice (1-6)"

$target = switch ($choice) {
    "1" { "esp32-mini32" }
    "2" { "xiao-esp32c3" }
    "3" { "pico" }
    "4" { "pico2" }
    "5" { "arduino-nano33" }
    "6" { Read-Host "Enter custom TinyGo board target name" }
    default { "esp32-mini32" }
}

Write-Host "`n[INFO] Compiling and flashing interpreter.go for target: $target..." -ForegroundColor Green

& tinygo flash "-target=$target" "interpreter.go"

if ($LASTEXITCODE -eq 0) {
    Write-Host "[SUCCESS] NanoSharp successfully flashed to $target!" -ForegroundColor Green
} else {
    Write-Host "[ERROR] Flashing failed. Check your USB connection and board drivers." -ForegroundColor Red
    exit 1
}

Write-Host "=== Flash Complete ===" -ForegroundColor Cyan