# goworld

*Not in active development! Working on it occasionally.*

Yet another attempt at building a basic OpenGL 3D engine from scratch, this time in Go. The goal is to create a engine capable of producing *some* kind of *passable* graphics with a coherent art style. So far, the style is based around colored voxels.

**Features:**
- Voxel editor demo with basic player physics
- Deferred Rendering Pipeline
- Directional Lights
- Directional Shadows
- Point Lights
- Screen-Space Ambient Occlusion (HBAO)
- Color Grading with Lookup Tables
- OBJ Model Loader
- TrueType Font Rendering
- React-like UI including a flexbox layout engine
- Custom ergonomic 3D math library derived from mathgl and go3d

Tested on OSX 10.10+ and Manjaro Linux. It should *theoretically* work on Windows.

![Screenshot from 2022-02-27](docs/img/screenshot220227.png)
![Screenshot from 2020-09-26](docs/img/screenshot200926.png)
![Screenshot from 2019-05-07](docs/img/screenshot190507.png)

## System Requirements

 * OpenGL 4.1 Core Profile

## Todo / Ideas

 * User Interface elements:
   * Textbox
 * Console
 * Lighting
   * Point Light shadows
   * Spot Light
   * Spot Light shadows
 * Bloom Post-Process Effect

## Building

### Open Dynamics Engine
- Download ODE 0.16.2
- Extract tarball
- Configure
  ```
  $ ./configure --enable-double-precision --enable-shared
  ```
- Compile
  ```
  make
  ```
- Install Shared Library
  ```
  make install
  ```
