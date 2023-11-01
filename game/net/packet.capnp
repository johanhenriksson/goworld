using Go = import "/go.capnp";
@0x9c4788c5a214e29c;
$Go.package("net");
$Go.import("github.com/johanhenriksson/goworld/game/net");

#
# Data types
#

struct Vec3 {
    x @0 :Float32;
    y @1 :Float32;
    z @2 :Float32;
}

#
# Packets
#

struct Packet {
    union {
        auth @0 :AuthPacket;
        move @1 :MovePacket;
    }
}

struct AuthPacket {
    token @0 :UInt64;
}

struct MovePacket {
    uid @0 :UInt64;
    position @1 :Vec3;
    rot @2 :Float32;
}
