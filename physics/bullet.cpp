#include "bullet.h"

#include "BulletCollision/CollisionDispatch/btGhostObject.h"
#include "BulletDynamics/Character/btKinematicCharacterController.h"
#include "_cgo_export.h"
#include "btBulletCollisionCommon.h"
#include "btBulletDynamicsCommon.h"
#include "bullet.hpp"

btVector3 vec3FromGo(goVector3 v) { return btVector3(v[0], v[1], v[2]); }

void vec3ToGo(btVector3& v, goVector3 out) {
    out[0] = v.getX();
    out[1] = v.getY();
    out[2] = v.getZ();
}

/* Dynamics World */
goDynamicsWorldHandle goCreateDynamicsWorld() {
    // these objects currently leak:
    auto colConfig = new btDefaultCollisionConfiguration();
    auto dispatcher = new btCollisionDispatcher(colConfig);
    auto broadphase = new btDbvtBroadphase();
    auto solver = new btSequentialImpulseConstraintSolver();

    auto world = new btDiscreteDynamicsWorld(dispatcher, broadphase, solver, colConfig);

    return (goDynamicsWorldHandle)world;
}

void goDeleteDynamicsWorld(goDynamicsWorldHandle world) {
    // todo: also clean up the other allocations, axisSweep,
    // pairCache,dispatcher,constraintSolver,collisionConfiguration
    btDynamicsWorld* dynamicsWorld = reinterpret_cast<btDynamicsWorld*>(world);
    delete dynamicsWorld;
}

void goSetGravity(goDynamicsWorldHandle world, goVector3 gravity) {
    btDynamicsWorld* dynamicsWorld = reinterpret_cast<btDynamicsWorld*>(world);
    btAssert(dynamicsWorld);
    dynamicsWorld->setGravity(btVector3(gravity[0], gravity[1], gravity[2]));
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

goRigidBodyHandle goCreateRigidBody(void* user_data, float mass, goShapeHandle cshape) {
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

    body->setUserPointer(user_data);

    return (goRigidBodyHandle)body;
}

void goDeleteRigidBody(goRigidBodyHandle cbody) {
    btRigidBody* body = reinterpret_cast<btRigidBody*>(cbody);
    btAssert(body);
    delete body;
}

/* Collision Shape definition */

goShapeHandle goNewSphereShape(goReal radius) { return (goShapeHandle) new btSphereShape(radius); }

goShapeHandle goNewBoxShape(goVector3 size) { return (goShapeHandle) new btBoxShape(vec3FromGo(size)); }

goShapeHandle goNewCylinderShape(goReal radius, goReal height) {
    return (goShapeHandle) new btCylinderShape(btVector3(radius, height, radius));
}

goShapeHandle goNewCapsuleShape(goReal radius, goReal height) {
    const int numSpheres = 2;
    btVector3 positions[numSpheres] = {btVector3(0, height, 0), btVector3(0, -height, 0)};
    btScalar radi[numSpheres] = {radius, radius};
    return (goShapeHandle) new btMultiSphereShape(positions, radi, numSpheres);
}

goShapeHandle goNewController(goVector3 position, goReal height, goReal stepHeight) {
    auto ghostObject = new btPairCachingGhostObject();
    btTransform transf;
    transf.setIdentity();
    ghostObject->setWorldTransform(transf);

    auto capsule = new btCapsuleShape(0.8, height);
    auto character = new btKinematicCharacterController(ghostObject, capsule, stepHeight);
    return (goShapeHandle)character;
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

goShapeHandle goNewStaticTriangleMeshShape(goTriangleMeshHandle meshHandle) {
    auto mesh = reinterpret_cast<btTriangleMesh*>(meshHandle);
    btAssert(mesh);
    return (goShapeHandle) new btBvhTriangleMeshShape(mesh, true);
}

goShapeHandle goNewCompoundShape() { return (goShapeHandle) new btCompoundShape(); }

void goAddChildShape(goShapeHandle compoundShapeHandle, goShapeHandle childShapeHandle, goVector3 childPos,
                     goQuaternion childOrn) {
    btCollisionShape* colShape = reinterpret_cast<btCollisionShape*>(compoundShapeHandle);
    btAssert(colShape->getShapeType() == COMPOUND_SHAPE_PROXYTYPE);

    btCompoundShape* compoundShape = reinterpret_cast<btCompoundShape*>(colShape);

    btCollisionShape* childShape = reinterpret_cast<btCollisionShape*>(childShapeHandle);

    btTransform localTrans;
    localTrans.setIdentity();
    localTrans.setOrigin(btVector3(childPos[0], childPos[1], childPos[2]));
    localTrans.setRotation(btQuaternion(childOrn[0], childOrn[1], childOrn[2], childOrn[3]));
    compoundShape->addChildShape(localTrans, childShape);
}

//	extern  void		goAddTriangle(goMeshInterfaceHandle meshHandle,
// goVector3 v0,goVector3 v1,goVector3 v2); 	extern  goShapeHandle
// goNewStaticTriangleMeshShape(goMeshInterfaceHandle);

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

void goSetScaling(goShapeHandle cshape, goVector3 cscaling) {
    btCollisionShape* shape = reinterpret_cast<btCollisionShape*>(cshape);
    btAssert(shape);

    btVector3 scaling(cscaling[0], cscaling[1], cscaling[2]);
    shape->setLocalScaling(scaling);
}

void goGetPosition(goRigidBodyHandle object, goVector3 position) {
    btRigidBody* body = reinterpret_cast<btRigidBody*>(object);
    btAssert(body);

    vec3ToGo(body->getWorldTransform().getOrigin(), position);
}

void goSetPosition(goRigidBodyHandle object, const goVector3 position) {
    btRigidBody* body = reinterpret_cast<btRigidBody*>(object);
    btAssert(body);

    btVector3 pos(position[0], position[1], position[2]);
    btTransform transf = body->getWorldTransform();
    transf.setOrigin(pos);
    body->setWorldTransform(transf);

    if (body->getMotionState()) {
        body->getMotionState()->setWorldTransform(transf);
    }

    body->activate();
}

void goGetRotation(goRigidBodyHandle body_handle, goVector3 rotation) {
    btRigidBody* body = reinterpret_cast<btRigidBody*>(body_handle);
    btAssert(body);

    btTransform transf = body->getWorldTransform();
    transf.getRotation().getEulerZYX(rotation[2], rotation[1], rotation[0]);
}

void goSetRotation(goRigidBodyHandle object, const goVector3 rotation) {
    btRigidBody* body = reinterpret_cast<btRigidBody*>(object);
    btAssert(body);

    btQuaternion orient;
    orient.setEulerZYX(rotation[2], rotation[1], rotation[0]);
    btTransform transf = body->getWorldTransform();
    transf.setRotation(orient);

    body->setWorldTransform(transf);
    if (body->getMotionState()) {
        body->getMotionState()->setWorldTransform(transf);
    }
    body->activate();
}

void goEnableDebug(goDynamicsWorldHandle world) {
    btDynamicsWorld* dynamicsWorld = reinterpret_cast<btDynamicsWorld*>(world);
    btAssert(dynamicsWorld);

    auto drawer = new GoDebugDrawer(GoDebugCallback);
    drawer->setDebugMode(btIDebugDraw::DBG_DrawWireframe | btIDebugDraw::DBG_DrawAabb);

    dynamicsWorld->setDebugDrawer(drawer);
}

void goDebugDraw(goDynamicsWorldHandle world) {
    btDynamicsWorld* dynamicsWorld = reinterpret_cast<btDynamicsWorld*>(world);
    btAssert(dynamicsWorld);
    dynamicsWorld->debugDrawWorld();
}
