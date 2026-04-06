package main

import (
	_ "embed"
	"fmt"
	"log"
	"math"
	"time"

	_ "github.com/Zyko0/go-sdl3/bin/binmix"
	_ "github.com/Zyko0/go-sdl3/bin/binsdl"
	_ "github.com/Zyko0/go-sdl3/bin/binttf"
	"github.com/Zyko0/go-sdl3/mixer"
	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/ttf"
)

//go:embed Lato-Bold.ttf
var latoBold []byte

//go:embed CantinaBand3.ogg
var soundData []byte

type JoystickDisplay struct {
	x, y       float32
	radius     float32
	joyX, joyY float32
	pressed    bool
}

func (j *JoystickDisplay) Render(renderer *sdl.Renderer) {
	col := sdl.FColor{R: 1, G: 0, B: 0, A: 1}
	if j.pressed {
		col = sdl.FColor{R: 1, G: 1, B: 0, A: 1}
	}
	renderFilledCircle(renderer, j.x, j.y, j.radius, col)

	joyposx := j.x + j.radius*j.joyX
	joyposy := j.y + j.radius*j.joyY
	renderFilledCircle(renderer, joyposx, joyposy, j.radius/5, sdl.FColor{R: 0, G: 1, B: 0, A: 1})
}

func renderFilledCircle(renderer *sdl.Renderer, cx, cy, r float32, col sdl.FColor) {
	const segments = 64
	vertices := make([]sdl.Vertex, segments+2)
	vertices[0] = sdl.Vertex{Position: sdl.FPoint{X: cx, Y: cy}, Color: col}
	for i := 0; i <= segments; i++ {
		angle := float64(i) / float64(segments) * 2 * math.Pi
		vertices[i+1] = sdl.Vertex{
			Position: sdl.FPoint{
				X: cx + r*float32(math.Cos(angle)),
				Y: cy + r*float32(math.Sin(angle)),
			},
			Color: col,
		}
	}
	indices := make([]int32, segments*3)
	for i := 0; i < segments; i++ {
		indices[i*3] = 0
		indices[i*3+1] = int32(i + 1)
		indices[i*3+2] = int32(i + 2)
	}
	renderer.RenderGeometry(nil, vertices, indices)
}

type TextDisplay struct {
	x, y    float32
	color   sdl.Color
	font    *ttf.Font
	text    string
	texture *sdl.Texture
	surface *sdl.Surface
}

func (t *TextDisplay) SetText(text string) error {
	t.Close()
	t.text = text
	surface, err := t.font.RenderTextBlended(t.text, t.color)
	if err != nil {
		log.Printf("Cannot set text: %s", err)
		return err
	}
	t.surface = surface
	return nil
}

func (t *TextDisplay) Render(renderer *sdl.Renderer) {
	if t.surface == nil {
		return
	}
	if t.texture == nil {
		texture, err := renderer.CreateTextureFromSurface(t.surface)
		if err != nil {
			log.Printf("cannot render text: %s", err)
			return
		}
		t.texture = texture
	}
	r := sdl.FRect{X: t.x, Y: t.y, W: float32(t.surface.W), H: float32(t.surface.H)}
	renderer.RenderTexture(t.texture, nil, &r)
}

func (t *TextDisplay) Close() {
	if t.texture != nil {
		t.texture.Destroy()
		t.texture = nil
	}
	if t.surface != nil {
		t.surface.Destroy()
		t.surface = nil
	}
}

func main() {
	if err := ttf.Init(); err != nil {
		log.Fatalf("TTF init: %s", err)
	}
	defer ttf.Quit()

	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS | sdl.INIT_AUDIO); err != nil {
		log.Fatalf("SDL init: %s", err)
	}
	defer sdl.Quit()

	window, renderer, err := sdl.CreateWindowAndRenderer("test", 640, 480, 0)
	if err != nil {
		log.Fatalf("CreateWindowAndRenderer: %s", err)
	}
	defer renderer.Destroy()
	defer window.Destroy()

	sdl.HideCursor()
	sdl.SetJoystickEventsEnabled(true)

	// Font
	fontStream, err := sdl.IOFromConstMem(latoBold)
	if err != nil {
		log.Fatalf("IOFromConstMem font: %s", err)
	}
	font, err := ttf.OpenFontIO(fontStream, true, 48)
	if err != nil {
		log.Fatalf("OpenFontIO: %s", err)
	}
	defer font.Close()

	text := TextDisplay{
		x:     10,
		y:     30,
		color: sdl.Color{R: 255, G: 0, B: 255, A: 255},
		font:  font,
	}

	// Audio
	var mixerDev *mixer.Mixer
	var musicTrack *mixer.Track

	if err := mixer.Init(); err != nil {
		log.Printf("Mixer init failed (audio disabled): %s", err)
	} else {
		defer mixer.Quit()

		mixerDev, err = mixer.CreateMixerDevice(sdl.AUDIO_DEVICE_DEFAULT_PLAYBACK, nil)
		if err != nil {
			log.Printf("CreateMixerDevice failed (audio disabled): %s", err)
		} else {
			defer mixerDev.Destroy()

			soundStream, err := sdl.IOFromConstMem(soundData)
			if err != nil {
				log.Printf("IOFromConstMem sound: %s", err)
			} else {
				audio, err := mixerDev.LoadAudio_IO(soundStream, true, true)
				if err != nil {
					log.Printf("LoadAudio_IO failed: %s", err)
				} else {
					defer audio.Destroy()

					musicTrack, err = mixerDev.CreateTrack()
					if err != nil {
						log.Printf("CreateTrack failed: %s", err)
						musicTrack = nil
					} else {
						defer musicTrack.Destroy()
						musicTrack.SetAudio(audio)
						musicTrack.SetGain(0.375) // ~48/128 from original SDL2 code
					}
				}
			}
		}
	}

	j1 := &JoystickDisplay{x: 100, y: 300, radius: 50}
	j2 := &JoystickDisplay{x: 540, y: 300, radius: 50}

	joysticks := make(map[sdl.JoystickID]*sdl.Joystick)

	running := true
	tick := time.Tick(time.Microsecond * 33333)

	for running {
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		j1.Render(renderer)
		j2.Render(renderer)
		text.Render(renderer)

		renderer.Present()

		<-tick

		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type {
			case sdl.EVENT_QUIT:
				fmt.Println("Quit")
				running = false
			case sdl.EVENT_JOYSTICK_AXIS_MOTION:
				ax := event.JoyAxisEvent()
				value := float32(ax.Value) / 32768.0
				switch ax.Axis {
				case 0:
					j1.joyX = value
				case 1:
					j1.joyY = value
				case 2:
					j2.joyX = value
				case 3:
					j2.joyY = value
				}
			case sdl.EVENT_JOYSTICK_BALL_MOTION:
				ball := event.JoyBallEvent()
				fmt.Println("Joystick", ball.Which, "trackball moved by", ball.Xrel, ball.Yrel)
			case sdl.EVENT_JOYSTICK_BUTTON_DOWN, sdl.EVENT_JOYSTICK_BUTTON_UP:
				jb := event.JoyButtonEvent()
				if jb.Down {
					text.SetText(fmt.Sprintf("Button %d/%d pressed", jb.Which, jb.Button))
				} else {
					text.SetText(fmt.Sprintf("Button %d/%d released", jb.Which, jb.Button))
				}
				if jb.Button == 14 {
					j1.pressed = jb.Down
				} else if jb.Button == 15 {
					j2.pressed = jb.Down
				} else if jb.Button == 2 && jb.Down && musicTrack != nil && !musicTrack.Playing() {
					musicTrack.Play(0)
				}
			case sdl.EVENT_JOYSTICK_ADDED:
				jd := event.JoyDeviceEvent()
				joystick, err := jd.Which.OpenJoystick()
				if err == nil {
					joysticks[jd.Which] = joystick
					fmt.Println("Joystick", jd.Which, "connected")
				}
			case sdl.EVENT_JOYSTICK_REMOVED:
				jd := event.JoyDeviceEvent()
				if joystick, ok := joysticks[jd.Which]; ok {
					joystick.Close()
					delete(joysticks, jd.Which)
				}
				fmt.Println("Joystick", jd.Which, "disconnected")
			}
		}
	}

	for _, joystick := range joysticks {
		joystick.Close()
	}
}
