using Go = import "/go.capnp";
@0x85d3acc39d94e0f8;
$Go.package("net");
$Go.import("github.com/johanhenriksson/goworld/game/net");

struct AuthPacket {
    token @0 :UInt64;
}

struct Vec3 {
    x @0 :Float32;
    y @1 :Float32;
    z @2 :Float32;
}

struct MovePacket {
    uid @0 :UInt64;
    position @1 :Vec3;
    rot @2 :Float32;
}
