# starroute

# run

```bash
starroute
```

# fullscreen

```bash
starroute -resize=full -both 1920x1080
```

# Minerals

Highest value minerals in the galaxy

| Mineral | Price per-gram |
|---------|----------------|
| Antihydrogen | 62T (trillion) |
| Technetium-99m | 2T (trillion) |
| Unbihexium-310 | 500B (billion) |
| Unbiquadium-304 | 100B (billion) |
| Unbiunium-299 | 50B (billion) |
| Oganesson-294 | 10B (billion) |
| Tennessine-294 | 5B (billion) |
| Livermorium-293 | 4B (billion) |
| Flerovium-289 | 3B (billion) |
| Copernicium-285 | 2.5B (billion) |
| Mendelevium-258 | 2B (billion) |
| Lawrencium-262 | 1.5B (billion) |
| Francium-223 | 1B (billion) |
| Nobelium-259 | 1B (billion) |
| Berkelium-247 | 350M |
| Californium-251 | 200M |
| Curium-248 | 175M |
| Fermium-257 | 150M |
| Astatine-211 | 125M |
| Astatine-210 | 100M |
| Plutonium-244 | 75M |
| Einsteinium-254 | 50M |
| Californium-250 | 45M |
| Californium-252 | 27M |
| Berkelium-249 | 25M |
| Red Diamond | 3M |
| Metastable Metallic Hydrogen | 2M |
| Einsteinium-253 | 1M |
| Neptunium-236 | 1M |
| Curium-246 | 700K |
| Scandium-47 | 450K |
| Curium-244 | 200K |
| Scandium-46 | 180K |
| Americium-243 | 160K |
| Helium-3 | 100K |
| Lutetium-177 | 80K |
| Tritium | 30K |
| Taaffeite | 20K |
| Promethium-145 | 15K |
| Plutonium-238 | 12K |
| Diamond | 10K |
| Uranium-233 | 8K |
| Neptunium-237 | 4K |
| Promethium-147 | 4K |
| Crystalline Osmium | 3K |
| Americium-241 | 2K |
| Rhodium | 500 |
| Scandium | 300 |
| Iridium | 150 |
| Osmium | 100 |
| Gold | 50 |
| Platinum | 40 |
| Palladium | 30 |
| Ruthenium | 20 |
| Rhenium | 15 |
| Silver | 2 |

# Ebiten References

- Hello world: https://ebitengine.org/en/tour/hello_world.html

- Tour: https://ebitengine.org/en/tour/

- Examples: https://ebitengine.org/en/examples/

- Cheatsheet: https://ebitengine.org/en/documents/cheatsheet.html

- Making games in Go: https://threedots.tech/post/making-games-in-go/

- Airplanes game: https://github.com/m110/airplanes

- Donburi ECS: https://github.com/yottahmd/donburi

- Tiled: https://www.mapeditor.org/

- Go Tiled Loader: https://github.com/lafriks/go-tiled

- Kenney Assets: https://kenney.nl/

- Ebiten UI: https://ebitenui.github.io

- Debug UI: https://github.com/ebitengine/debugui

- Short sound effect

Creating an audio.Player is not expensive. It is fine to create one player for one short sound effect. For example, this code is totally fine:

```golang
// PlaySE plays a sound effect.
func PlaySE(bs []byte) {
    sePlayer := audioContext.NewPlayerFromBytes(bs)
    // sePlayer is never GCed as long as it plays.
    sePlayer.Play()
}
```

https://ebitengine.org/en/documents/performancetips.html
