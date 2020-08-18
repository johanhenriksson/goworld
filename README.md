# goworld

*Not in active development! Working on it occasionally.*

Yet another attempt at building a basic OpenGL 3D engine from scratch, Google Go edition. The goal is to create a engine capable of producing *some* kind of *passable* graphics with a coherent art style. So far, the style is based on colored voxels.

**Features:**
- Deferred Rendering Pipeline
- Directional Lights
- Directional Shadows
- Point Lights
- Screen-Space Ambient Occlusion (HBAO)
- Color Grading with Lookup Tables
- OBJ Model Loader
- TrueType Font Rendering
- UI: Panels, Labels, basic layout engine

Only tested on OSX 10.10+. It should theoretically work on Windows/Linux.

![Screenshot from 2019-05/07](docs/img/screenshot190507.png)

## System Requirements

 * OpenGL 4.1 Core Profile

## Todo / Ideas

 * Improved scene graph
 * Components System
 * User Interface elements:
   * Textbox
   * Button
 * Embedded scripting language (probably V8)
 * Console
 * Lighting
   * Point Light shadows
   * Spot Light
   * Spot Light shadows
 * Bloom Effect
 * Scene save/load
