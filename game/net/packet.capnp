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
        unknown @0: Void;

        auth @1 :AuthPacket;

        entityMove @2 :EntityMovePacket;
        entitySpawn @3 :EntitySpawnPacket;
        entityObserve @4 :EntityObservePacket;
    }
}

struct AuthPacket {
    token @0 :UInt64;
}

struct EntityMovePacket {
    entity @0 :UInt64;
    position @1 :Vec3;
    rotation @2 :Float32;
    stopped @3 :Bool;
}

struct EntitySpawnPacket {
    entity @0 :UInt64;
    position @1 :Vec3;
}

struct EntityObservePacket {
    entity @0 :UInt64;
}
