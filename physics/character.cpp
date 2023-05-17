#include "BulletCollision/CollisionDispatch/btGhostObject.h"
#include "BulletDynamics/Character/btKinematicCharacterController.h"
#include "btBulletCollisionCommon.h"
#include "btBulletDynamicsCommon.h"
#include "bullet.h"

goCharacterHandle goCreateCharacter(goShapeHandle shapeHandle, float height, float radius, float stepHeight) {
    btConvexShape* capsule = new btCapsuleShape(radius, height);
    btAssert(capsule->getShapeType() == CONVEX_HULL_SHAPE_PROXYTYPE);

    btTransform trans;
    trans.setIdentity();
    trans.setOrigin(btVector3(5, 15, 5));

    auto ghostObject = new btPairCachingGhostObject();
    ghostObject->setWorldTransform(trans);
    ghostObject->setCollisionShape(capsule);
    ghostObject->setCollisionFlags(btCollisionObject::CF_CHARACTER_OBJECT);

    btKinematicCharacterController* character =
        new btKinematicCharacterController(ghostObject, capsule, stepHeight, btVector3(0, 1, 0));

    return (goCharacterHandle)character;
}

void goDeleteCharacter(goCharacterHandle handle) {
    btKinematicCharacterController* character = reinterpret_cast<btKinematicCharacterController*>(handle);
    btAssert(character);

    delete character->getGhostObject();
    delete character;
}

void goAddCharacter(goDynamicsWorldHandle worldHandle, goCharacterHandle handle) {
    btDynamicsWorld* world = reinterpret_cast<btDynamicsWorld*>(worldHandle);
    btAssert(world);
    btKinematicCharacterController* character = reinterpret_cast<btKinematicCharacterController*>(handle);
    btAssert(character);
    btAssert(character->getGhostObject());

    world->addCollisionObject(character->getGhostObject(), btBroadphaseProxy::CharacterFilter,
                              btBroadphaseProxy::StaticFilter | btBroadphaseProxy::DefaultFilter);

    world->addAction(character);
}

void goRemoveCharacter(goDynamicsWorldHandle worldHandle, goCharacterHandle handle) {
    btDynamicsWorld* world = reinterpret_cast<btDynamicsWorld*>(worldHandle);
    btAssert(world);
    btKinematicCharacterController* character = reinterpret_cast<btKinematicCharacterController*>(handle);
    btAssert(character);

    // detach from world
    world->removeAction(character);

    // clean up ghost object
    world->removeCollisionObject(character->getGhostObject());
}

void goCharacterWalkDirection(goCharacterHandle handle, goVector3 direction) {
    btKinematicCharacterController* character = reinterpret_cast<btKinematicCharacterController*>(handle);
    btAssert(character);
    character->setWalkDirection(btVector3(direction[0], direction[1], direction[2]));
}

void goCharacterJump(goCharacterHandle handle) {
    btKinematicCharacterController* character = reinterpret_cast<btKinematicCharacterController*>(handle);
    btAssert(character);
    character->jump();
}

void goCharacterWarp(goCharacterHandle handle, goVector3 position) {
    btKinematicCharacterController* character = reinterpret_cast<btKinematicCharacterController*>(handle);
    btAssert(character);
    character->warp(btVector3(position[0], position[1], position[2]));
}

void goCharacterUpdate(goCharacterHandle handle, goDynamicsWorldHandle worldHandle, float dt) {
    btKinematicCharacterController* character = reinterpret_cast<btKinematicCharacterController*>(handle);
    btAssert(character);
    auto world = reinterpret_cast<btDynamicsWorld*>(worldHandle);
    btAssert(world);
    character->playerStep(world, dt);
}
