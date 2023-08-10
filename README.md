# goworld

An attempt at building a basic 3D engine from scratch in Go using Vulkan.

**Goals**

- Have fun making it
- Use minimal dependencies
- Implement modern GPU-driven rendering techniques
- Ideally leverage Go's concurrency useful ways
- Run on MacOS/Linux/Windows
- Create an ergonomic, reactive GUI system
- Get to a state capable of producing _some_ kind of _passable_ graphics with a coherent art style
- Experiment with some cool demo scenes
- Maybe make a game

**Features**

Currently, the following features exist in varying degrees of completeness:

- Classic object/component scene graph
- Basic scene editor
- Voxel world/editor demo
- Unified Rendering Pipeline (Forward/Deferred)
  - Directional Lights (w/ cascading shadow maps)
  - Point Lights
  - Screen-Space Ambient Occlusion (HBAO)
  - Color Grading with Lookup Tables
- 3D Physics Engine integration via Bullet SDK
  - Character controller
  - Rigidbody dynamics
  - Basic shape colliders (box, sphere, cylinder, capsule)
  - Mesh colliders
- TrueType Font Rendering
- React-like UI with a flexbox layout engine & css-like styling
  - Hooks
  - Portals/fragments
  - Rect
  - Image
  - Button
  - Textbox
  - Floating windows
- Custom ergonomic 3D math library derived from mathgl and go3d

![Screenshot from 2023-02-06](docs/img/screenshot230305.png)
![Screenshot from 2022-02-27](docs/img/screenshot220227.png)
![Screenshot from 2020-09-26](docs/img/screenshot200926.png)
![Screenshot from 2019-05-07](docs/img/screenshot190507.png)

## System Requirements

- Vulkan 1.2
- MacOS users need MoltenVK.

## Build Instructions

Goworld is developed & tested on MacOS 13. It should be reasonably easy to get it running on Linux or Windows,
but its not officially supported yet.

### (MacOS) Install MoltenVK

Grab the latest version of MoltenVK.

### Build Bullet SDK

Goworld uses the Bullet SDK for physics. In order to compile from source, you first need to compile Bullet.

```bash
# check out bullet3
git clone https://github.com/bulletphysics/bullet3.git
cd bullet3

# configure build
cmake . \
    -DBUILD_SHARED_LIBS=ON \
    -DINSTALL_LIBS=ON \
    -DUSE_DOUBLE_PRECISION=OFF \
    -DBUILD_BULLET2_DEMOS=OFF \
    -DBUILD_CPU_DEMOS=OFF \
    -DBUILD_OPENGL3_DEMOS=OFF \
    -DBUILD_BULLET3=ON \
    -DBUILD_PYBULLET=OFF \
    -DBUILD_EXTRAS=OFF \
    -G "Unix Makefiles"

# compile & install
make
make install
```

### Build

Goworld uses Taskfile for convenient building:

- Build:
  ```
  $ task build
  ```
- Build & run:
  ```
  $ task run
  ```
