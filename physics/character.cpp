#include <stdbool.h>

#include "BulletCollision/CollisionDispatch/btGhostObject.h"
#include "BulletDynamics/Character/btKinematicCharacterController.h"
#include "btBulletCollisionCommon.h"
#include "btBulletDynamicsCommon.h"
#include "bullet.h"

goCharacterHandle goCreateCharacter(goShapeHandle shapeHandle, float stepHeight) {
    auto shape = reinterpret_cast<btConvexShape*>(shapeHandle);

    auto ghostObject = new btPairCachingGhostObject();
    ghostObject->setCollisionShape(shape);
    ghostObject->setCollisionFlags(btCollisionObject::CF_CHARACTER_OBJECT);

    auto character = new btKinematicCharacterController(ghostObject, shape, stepHeight, btVector3(0, 1, 0));

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
                              btBroadphaseProxy::AllFilter);

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

void goCharacterMove(goCharacterHandle handle, goVector3* direction) {
    btKinematicCharacterController* character = reinterpret_cast<btKinematicCharacterController*>(handle);
    btAssert(character);
    character->setWalkDirection(btVector3(direction->x, direction->y, direction->z));
}

void goCharacterJump(goCharacterHandle handle) {
    btKinematicCharacterController* character = reinterpret_cast<btKinematicCharacterController*>(handle);
    btAssert(character);
    character->jump();
}

void goCharacterGetState(goCharacterHandle handle, goCharacterState* state) {
    btKinematicCharacterController* character = reinterpret_cast<btKinematicCharacterController*>(handle);
    btAssert(character);

    auto transf = character->getGhostObject()->getWorldTransform();

    auto pos = transf.getOrigin();
    state->position.x = pos.x();
    state->position.y = pos.y();
    state->position.z = pos.z();

    auto rot = transf.getRotation();
    state->rotation.x = rot.x();
    state->rotation.y = rot.y();
    state->rotation.z = rot.z();
    state->rotation.w = rot.w();

    state->grounded = character->onGround();
}

void goCharacterSetState(goCharacterHandle handle, goCharacterState* state) {
    btKinematicCharacterController* character = reinterpret_cast<btKinematicCharacterController*>(handle);
    btAssert(character);

    auto transf = character->getGhostObject()->getWorldTransform();

    auto rot = btQuaternion(state->rotation.x, state->rotation.y, state->rotation.z, state->rotation.w);
    transf.setRotation(rot);

    auto pos = btVector3(state->position.x, state->position.y, state->position.z);
    transf.setOrigin(pos);

    character->getGhostObject()->setWorldTransform(transf);
}
