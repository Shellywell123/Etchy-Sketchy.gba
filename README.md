# Etchy-Sketchy

Homebrew Game boy Advance game written in TinyGo. 

<p float="left">
  <img src="https://github.com/Shellywell123/Etchy-Sketchy/blob/main/assets/ScreenGrab.GIF" width="400" />
  <img src="assets/GameBoyAdvance.png" width="400" />
</p>

## Play 
- Play online (Coming Soon)\
or
- Download the latest [ROM](https://github.com/Shellywell123/Etchy-Sketchy/releases) and play it in an emulator/original hardware.

## Build ROM
```bash
cd src && tinygo build -size short -o Ethcy-Sketchy.gba -target=gameboy-advance etchy_sketchy.go
```

## References
Quick thanks to the other community projects that made this project so simple.
- [scraly/learning-go-by-examples](https://github.com/scraly/learning-go-by-examples) for the example template.
- [danmrichards/go-gba-pong](https://github.com/danmrichards/go-gba-pong) for the key press method.
