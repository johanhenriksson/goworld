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
