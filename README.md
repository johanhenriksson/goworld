# goworld

*Not in active development! Working on it occasionally.*

Yet another attempt at building a basic OpenGL 3D engine from scratch, this time in Go. The goal is to create a engine capable of producing *some* kind of *passable* graphics with a coherent art style. So far, the style is based around colored voxels.

**Features:**
- Voxel World demo, with basic player physics and an editable, persistent world.
- Deferred Rendering Pipeline
- Directional Lights
- Directional Shadows
- Point Lights
- Screen-Space Ambient Occlusion (HBAO)
- Color Grading with Lookup Tables
- OBJ Model Loader
- TrueType Font Rendering
- UI: Panels, Labels, Images, and a simple layout engine
- Custom ergonomic 3D math library derived from mathgl and go3d

Only tested on OSX 10.10+. It should *theoretically* work on Windows and Linux.

![Screenshot from 2020-09-26](docs/img/screenshot200926.png)
![Screenshot from 2019-05-07](docs/img/screenshot190507.png)

## System Requirements

 * OpenGL 4.1 Core Profile

## Todo / Ideas

 * Improved Scene Graph / Components System
 * User Interface elements:
   * Textbox
   * Button
 * Embedded scripting language (probably javascript/V8)
 * Console
 * Lighting
   * Point Light shadows
   * Spot Light
   * Spot Light shadows
 * Bloom Post-Process Effect
