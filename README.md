# goworld

*Recently revived after about 3 years*

Attempt at building a basic OpenGL 3D engine from scratch, Google Go edition. The goal is to create a engine capable of producing *some* kind of *passable* graphics with a coherent art style. So far, the style has been based on colored voxels.

Only tested on OSX 10.10+. It should theoretically work on Windows/Linux.

![Screenshot from 2019-05/07](docs/img/screenshot190507.png)

## System Requirements

 * OpenGL 4.1 Core Profile

## Dependencies

 * GLFW
 * OpenGL 4.1
 * Open Dynamics Engine (ODE)

**Installation Steps**
 
 * Clone & build ODE
 * Install dependencies
 * Run!

## Todo / Ideas

 * ~~Basic Scene graph~~
 * ~~Defered rendering pipeline~~
 * ~~Basic resource management:~~
   * ~~Shaders~~
   * ~~Textures~~
   * ~~Materials~~
 * Improved scene graph
 * Components System
 * User Interface elements:
   * ~~Frame~~
   * ~~Text Label~~
   * ~~Image~~
   * Textbox
   * Button
   * Keyboard events
   * Mouse events
 * Embedded language (probably V8)
 * Console
 * Lighting
   * ~~Directional light~~
   * ~~Directional light shadows~~
   * ~~Point light~~
   * Point light shadows
   * ~~Screen-Space Ambient Occlusion~~
 * Bloom Effect
 * Scene save/load
