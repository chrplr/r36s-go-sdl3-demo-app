# R36S SDL3 Demo App

This is a port to go-sdl3 of the code at <https://github.com/AndreRenaud/r36s-demo-app> demonstrating the feasibility of running go-sdlX programs to the r36s. See <https://ignavus.net/r36s> for more info about the original project.

The program is a joystick input demo for the R36S handheld, built with [go-sdl3](https://github.com/Zyko0/go-sdl3). Displays two analog stick visualizations, reacts to button presses with text feedback, and plays audio. SDL3 is bundled inside the binary — no system SDL3 installation needed on the device.

> As of 2026-04-06, I manage to compile and run the app on the R36S, but there is no reaction to the joystick input, so the program hangs, displaying two circles. The same thing happens on my linux laptop with a joystick plugged in, so I suppose the issue is not specific to the R36S, but rather to the port (done essentially by Claude).

## Prerequisites (build machine)

- Go 1.22+

## Building

```sh
GOOS=linux GOARCH=arm64 go build -o r36s-demo-app .
```

No CGo, no cross-compiler toolchain needed.

## Installing on a R36S running ArkOS4Clone

To make the app available on the R36S. 

1. Insert the SD card on your PC and set the variable EASYROMS to point to the EASYROM partition.

   ```sh
   EASYROMS=/media/xxxxxxx/EASYROMS
   ``` 

2. Create a folder for the app in  `ports`, e.g. `myapp` and copy the binary there:

   ```sh
   mkdir -p $EASYROMS/ports/myapp
   cp r36-demo-app $EASYROMS/ports/myapp
   chmod +x $EASYROMS/ports/myapp/r36-demo-app
   ```
 
3. Create a launcher script in `ports` (not in `myapp`!)

   ```sh
   cat > $EASYROMS/ports/r36s-demo-app.sh << 'EOF'
   #!/bin/bash
   /roms/ports/myapp/r36s-demo-app > /roms/ports/myapp/errors.log 2>&1
   EOF
   chmod +x $EASYROMS/r36s-demo-app.sh
   ```

4. Put the SD card back in the r36s, and among the emulators, select "Ports". You should be able to find the app there.
   If it does not work, check `errors.log` in `ports/myapp`


## Controls

| Input | Action |
|-------|--------|
| Left stick | Move left joystick indicator |
| Right stick | Move right joystick indicator |
| Button 14 | Highlight left joystick (yellow) |
| Button 15 | Highlight right joystick (yellow) |
| Button 2 (B) | Play Cantina Band (if not already playing) |
