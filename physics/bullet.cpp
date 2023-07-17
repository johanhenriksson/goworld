#include "bullet.h"

#include <iostream>
#include <ostream>

#include "BulletCollision/CollisionDispatch/btGhostObject.h"
#include "BulletDynamics/Character/btKinematicCharacterController.h"
#include "_cgo_export.h"
#include "btBulletCollisionCommon.h"
#include "btBulletDynamicsCommon.h"
#include "bullet.hpp"

btVector3 vec3FromGo(goVector3* v) { return btVector3(v->x, v->y, v->z); }

void vec3ToGo(btVector3& v, goVector3* out) {
    out->x = v.getX();
    out->y = v.getY();
    out->z = v.getZ();
}

btQuaternion quatFromGo(goQuaternion* q) { return btQuaternion(q->x, q->y, q->z, q->w); }

void quatToGo(btQuaternion& q, goQuaternion* out) {
    out->x = q.getX();
    out->y = q.getY();
    out->z = q.getZ();
    out->w = q.getW();
}

/* Dynamics World */
goDynamicsWorldHandle goCreateDynamicsWorld() {
    // these objects currently leak:
    auto colConfig = new btDefaultCollisionConfiguration();
    auto dispatcher = new btCollisionDispatcher(colConfig);
    auto solver = new btSequentialImpulseConstraintSolver();

    auto broadphase = new btDbvtBroadphase();

    // this line is required to use the character controller
    broadphase->getOverlappingPairCache()->setInternalGhostPairCallback(new btGhostPairCallback());

    auto world = new btDiscreteDynamicsWorld(dispatcher, broadphase, solver, colConfig);

    return (goDynamicsWorldHandle)world;
}

void goDeleteDynamicsWorld(goDynamicsWorldHandle world) {
    // todo: also clean up the other allocations, axisSweep,
    // pairCache,dispatcher,constraintSolver,collisionConfiguration
    btDynamicsWorld* dynamicsWorld = reinterpret_cast<btDynamicsWorld*>(world);
    delete dynamicsWorld;
}

void goSetGravity(goDynamicsWorldHandle world, goVector3* gravity) {
    btDynamicsWorld* dynamicsWorld = reinterpret_cast<btDynamicsWorld*>(world);
    btAssert(dynamicsWorld);
    dynamicsWorld->setGravity(vec3FromGo(gravity));
}

void goStepSimulation(goDynamicsWorldHandle world, goReal timeStep) {
    btDynamicsWorld* dynamicsWorld = reinterpret_cast<btDynamicsWorld*>(world);
    btAssert(dynamicsWorld);
    dynamicsWorld->stepSimulation(timeStep);
}

void goAddRigidBody(goDynamicsWorldHandle world, goRigidBodyHandle object) {
    btDynamicsWorld* dynamicsWorld = reinterpret_cast<btDynamicsWorld*>(world);
    btAssert(dynamicsWorld);
    btRigidBody* body = reinterpret_cast<btRigidBody*>(object);
    btAssert(body);

    dynamicsWorld->addRigidBody(body);
}

void goRemoveRigidBody(goDynamicsWorldHandle world, goRigidBodyHandle object) {
    btDynamicsWorld* dynamicsWorld = reinterpret_cast<btDynamicsWorld*>(world);
    btAssert(dynamicsWorld);
    btRigidBody* body = reinterpret_cast<btRigidBody*>(object);
    btAssert(body);

    dynamicsWorld->removeRigidBody(body);
}

/* Rigid Body  */

goRigidBodyHandle goCreateRigidBody(char* user, float mass, goShapeHandle cshape) {
    btCollisionShape* shape = reinterpret_cast<btCollisionShape*>(cshape);
    btAssert(shape);

    btTransform trans;
    trans.setIdentity();

    btVector3 localInertia(0, 0, 0);
    if (mass > 0) {
        // calculate inertia if dynamic (mass > 0)
        shape->calculateLocalInertia(mass, localInertia);
    }

    auto motionState = new btDefaultMotionState(trans);

    btRigidBody::btRigidBodyConstructionInfo rbci(mass, motionState, shape, localInertia);
    btRigidBody* body = new btRigidBody(rbci);

    body->setUserPointer(user);

    return (goRigidBodyHandle)body;
}

void goDeleteRigidBody(goRigidBodyHandle cbody) {
    btRigidBody* body = reinterpret_cast<btRigidBody*>(cbody);
    btAssert(body);
    delete body;
}

void goRigidBodyGetState(goRigidBodyHandle object, goRigidBodyState* state) {
    btRigidBody* body = reinterpret_cast<btRigidBody*>(object);
    btAssert(body);

    auto transf = body->getWorldTransform();
    vec3ToGo(transf.getOrigin(), &state->position);
    auto rot = transf.getRotation();
    quatToGo(rot, &state->rotation);
}

void goRigidBodySetShape(goRigidBodyHandle objectPtr, goShapeHandle shapePtr) {
    btRigidBody* body = reinterpret_cast<btRigidBody*>(objectPtr);
    btAssert(body);

    btCollisionShape* shape = reinterpret_cast<btCollisionShape*>(shapePtr);
    btAssert(shape);

    body->setCollisionShape(shape);
}

void goRigidBodySetState(goRigidBodyHandle object, goRigidBodyState* state) {
    btRigidBody* body = reinterpret_cast<btRigidBody*>(object);
    btAssert(body);

    auto transf = body->getWorldTransform();
    auto pos = vec3FromGo(&state->position);
    auto rot = quatFromGo(&state->rotation);
    transf.setOrigin(pos);
    transf.setRotation(rot);
    body->setWorldTransform(transf);
}

/* Collision Shape definition */

goShapeHandle goNewSphereShape(char* user, goReal radius) {
    auto shape = new btSphereShape(radius);
    shape->setUserPointer(user);
    return (goShapeHandle)shape;
}

goShapeHandle goNewBoxShape(char* user, goVector3* size) {
    auto shape = new btBoxShape(vec3FromGo(size));
    shape->setUserPointer(user);
    return (goShapeHandle)shape;
}

goShapeHandle goNewCylinderShape(char* user, goReal radius, goReal height) {
    auto shape = new btCylinderShape(btVector3(radius, height, radius));
    shape->setUserPointer(user);
    return (goShapeHandle)shape;
}

goShapeHandle goNewCapsuleShape(char* user, goReal radius, goReal height) {
    auto shape = new btCapsuleShape(radius, height);
    shape->setUserPointer(user);
    return (goShapeHandle)shape;
}

/* Concave static triangle meshes */
goTriangleMeshHandle goNewTriangleMesh(void* vertex_ptr, int vertex_count, int vertex_stride, void* index_ptr,
                                       int index_count, int index_stride) {
    auto mesh = btIndexedMesh();

    mesh.m_numVertices = vertex_count;
    mesh.m_vertexBase = (const unsigned char*)vertex_ptr;
    mesh.m_vertexStride = vertex_stride;

    mesh.m_numTriangles = index_count / 3;
    mesh.m_triangleIndexBase = (const unsigned char*)index_ptr;
    mesh.m_triangleIndexStride = 3 * index_stride;

    // infer index type from its stride. this is kinda stupid
    PHY_ScalarType indexType = PHY_INTEGER;
    switch (index_stride) {
        case 1:
            indexType = PHY_UCHAR;
        case 2:
            indexType = PHY_SHORT;
    }

    auto array = new btTriangleIndexVertexArray();
    array->addIndexedMesh(mesh, indexType);
    return (goTriangleMeshHandle)array;
}

goShapeHandle goNewTriangleMeshShape(char* user, goTriangleMeshHandle meshHandle) {
    auto mesh = reinterpret_cast<btTriangleMesh*>(meshHandle);
    btAssert(mesh);
    auto shape = new btBvhTriangleMeshShape(mesh, true);
    shape->setUserPointer(shape);
    return (goShapeHandle)shape;
}

goShapeHandle goNewCompoundShape() { return (goShapeHandle) new btCompoundShape(); }

void goAddChildShape(goShapeHandle compoundShapeHandle, goShapeHandle childShapeHandle, goVector3* childPos,
                     goQuaternion* childOrn) {
    btCollisionShape* colShape = reinterpret_cast<btCollisionShape*>(compoundShapeHandle);
    btAssert(colShape->getShapeType() == COMPOUND_SHAPE_PROXYTYPE);

    btCompoundShape* compoundShape = reinterpret_cast<btCompoundShape*>(colShape);

    btCollisionShape* childShape = reinterpret_cast<btCollisionShape*>(childShapeHandle);

    btTransform localTrans;
    localTrans.setIdentity();
    localTrans.setOrigin(vec3FromGo(childPos));
    localTrans.setRotation(quatFromGo(childOrn));
    compoundShape->addChildShape(localTrans, childShape);
}

void goAddVertex(goShapeHandle cshape, goReal x, goReal y, goReal z) {
    btCollisionShape* colShape = reinterpret_cast<btCollisionShape*>(cshape);
    btAssert(colShape->getShapeType() == CONVEX_HULL_SHAPE_PROXYTYPE);

    btConvexHullShape* convexHullShape = reinterpret_cast<btConvexHullShape*>(cshape);
    convexHullShape->addPoint(btVector3(x, y, z));
}

void goDeleteShape(goShapeHandle cshape) {
    btCollisionShape* shape = reinterpret_cast<btCollisionShape*>(cshape);
    btAssert(shape);
    delete shape;
}

void goSetScaling(goShapeHandle cshape, goVector3* cscaling) {
    btCollisionShape* shape = reinterpret_cast<btCollisionShape*>(cshape);
    btAssert(shape);

    shape->setLocalScaling(vec3FromGo(cscaling));
}

void goEnableDebug(goDynamicsWorldHandle world) {
    btDynamicsWorld* dynamicsWorld = reinterpret_cast<btDynamicsWorld*>(world);
    btAssert(dynamicsWorld);

    auto drawer = dynamicsWorld->getDebugDrawer();
    if (!drawer) {
        auto drawline = [world](const btVector3& from, const btVector3 to, const btVector3 color) -> void {
            GoDrawLineCallback(world, from.getX(), from.getY(), from.getZ(), to.getX(), to.getY(), to.getZ(),
                               color.getX(), color.getY(), color.getZ());
        };

        drawer = new GoDebugDrawer(drawline);
        dynamicsWorld->setDebugDrawer(drawer);
    }

    drawer->setDebugMode(btIDebugDraw::DBG_DrawWireframe | btIDebugDraw::DBG_DrawAabb);
}

void goDisableDebug(goDynamicsWorldHandle world) {
    btDynamicsWorld* dynamicsWorld = reinterpret_cast<btDynamicsWorld*>(world);
    btAssert(dynamicsWorld);

    auto drawer = dynamicsWorld->getDebugDrawer();
    if (drawer) {
        drawer->setDebugMode(btIDebugDraw::DBG_NoDebug);
    }
}

void goDebugDraw(goDynamicsWorldHandle world) {
    btDynamicsWorld* dynamicsWorld = reinterpret_cast<btDynamicsWorld*>(world);
    btAssert(dynamicsWorld);
    dynamicsWorld->debugDrawWorld();
}

int goRayCast(goDynamicsWorldHandle world, goVector3* rayStart, goVector3* rayEnd, goRayCastResult* res) {
    btDynamicsWorld* dynamicsWorld = reinterpret_cast<btDynamicsWorld*>(world);
    btAssert(dynamicsWorld);

    auto from = vec3FromGo(rayStart);
    auto to = vec3FromGo(rayEnd);
    auto callback = btCollisionWorld::ClosestRayResultCallback(from, to);

    dynamicsWorld->rayTest(from, to, callback);

    if (callback.hasHit()) {
        vec3ToGo(callback.m_hitPointWorld, &res->point);
        vec3ToGo(callback.m_hitNormalWorld, &res->normal);
        auto shape = callback.m_collisionObject->getCollisionShape();
        res->shape = shape->getUserPointer();
        return 1;
    }

    return 0;
}
