# goworld

Yet another attempt at building a basic 3D engine from scratch, this time in Go. The goal is to create a engine capable of producing _some_ kind of _passable_ graphics with a coherent art style. So far, the style is based around colored voxels.

**Features:**

- Voxel world/editor demo
- Deferred Rendering Pipeline
- 3D Physics Engine via Bullet SDK
- Directional Lights
- Directional Shadows
- Point Lights
- TrueType Font Rendering
- React-like UI including a flexbox layout engine
- Custom ergonomic 3D math library derived from mathgl and go3d
- Screen-Space Ambient Occlusion (HBAO)
- Color Grading with Lookup Tables

Tested on OSX 10.10+ and Manjaro Linux. It should _theoretically_ work on Windows.

![Screenshot from 2023-02-06](docs/img/screenshot230305.png)
![Screenshot from 2022-02-27](docs/img/screenshot220227.png)
![Screenshot from 2020-09-26](docs/img/screenshot200926.png)
![Screenshot from 2019-05-07](docs/img/screenshot190507.png)

## System Requirements

- Vulkan 1.2
- MacOS users need MoltenVK.

## Build Bullet SDK

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
