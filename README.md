```
████████╗███████╗██████╗ ███╗   ███╗██████╗ ███████╗███╗   ███╗ ██████╗
╚══██╔══╝██╔════╝██╔══██╗████╗ ████║██╔══██╗██╔════╝████╗ ████║██╔═══██╗
   ██║   █████╗  ██████╔╝██╔████╔██║██║  ██║█████╗  ██╔████╔██║██║   ██║
   ██║   ██╔══╝  ██╔══██╗██║╚██╔╝██║██║  ██║██╔══╝  ██║╚██╔╝██║██║   ██║
   ██║   ███████╗██║  ██║██║ ╚═╝ ██║██████╔╝███████╗██║ ╚═╝ ██║╚██████╔╝
   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝╚═════╝ ╚══════╝╚═╝     ╚═╝ ╚═════╝
```
# Termdemo: 80x24 24bit 60fps
This repository holds collections of demos for Linux terminals. The motivation
is to encourage exploring everything possible under `xterm-ghostty`, meaning
the use of glyphs, unicode, font bold/italics, full 24-bit colors, and  pushing
for 60fps performance at a classic 80x24 resolution. Other terminals may not be
fully compatible. `go run ./modernretro` to see the showcase demo.

**Technical Constraints:**
- Targeted for `TERM=xterm-256colors` or `TERM=xterm-ghostty`.
- Designed for an 80x24 terminal size with 24-bit color and 60fps capability.
- No kitty or sixel graphics. Different kind of demo

Code here is not important, just showcasing that is possible under the above
restrictions, also moving away just "AAlib" style rasteration of images into
ASCII, but explore the use of unicode characters, font styles, and color to
create more dynamic and visually appealing demos.
