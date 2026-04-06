# R36S SDL3 Demo App

A joystick input demo for the R36S handheld, built with [go-sdl3](https://github.com/Zyko0/go-sdl3). Displays two analog stick visualizations, reacts to button presses with text feedback, and plays audio. SDL3 is bundled inside the binary — no system SDL3 installation needed on the device.

## Prerequisites (build machine)

- Go 1.22+
- `sshpass` (for the deploy script): `apt install sshpass` / `brew install sshpass`

## Building

```sh
GOOS=linux GOARCH=arm64 go build -o r36s-demo-app .
```

No CGo, no cross-compiler toolchain needed.

## Installing on a R36S running ArkOS4Clone

### 1. Connect to WiFi

On the device, go to **Options → WiFi** and connect to your network. Note the IP address shown, or find it from your router's DHCP table.

### 2. Enable SSH

SSH is enabled by default on ArkOS4Clone. The credentials are:

- **User**: `root`
- **Password**: `root`

### 3. Copy and run the binary

Use the deploy script (requires `sshpass`):

```sh
./original-sdl2/run_on_unit.sh <device-ip>
```

This builds the ARM64 binary, kills any running instance, copies it to `/tmp/` on the device, and launches it over SSH.

Or do it manually:

```sh
scp r36s-demo-app root@<device-ip>:/tmp/
ssh root@<device-ip> "cd /tmp && ./r36s-demo-app"
```

### 4. Permanent installation

To make the app available as a port in EmulationStation, copy the binary to the ports directory and create a launch script:

```sh
ssh root@<device-ip>
cp /tmp/r36s-demo-app /roms/ports/
cat > /roms/ports/r36s-demo-app.sh << 'EOF'
#!/bin/bash
/roms/ports/r36s-demo-app
EOF
chmod +x /roms/ports/r36s-demo-app.sh
```

Then refresh the Ports section in EmulationStation (or restart it).

## Controls

| Input | Action |
|-------|--------|
| Left stick | Move left joystick indicator |
| Right stick | Move right joystick indicator |
| Button 14 | Highlight left joystick (yellow) |
| Button 15 | Highlight right joystick (yellow) |
| Button 2 (B) | Play Cantina Band (if not already playing) |
