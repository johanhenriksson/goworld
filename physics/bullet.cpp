
#include "bullet.h"

#include "BulletCollision/CollisionDispatch/btGhostObject.h"
#include "BulletDynamics/Character/btKinematicCharacterController.h"
#include "btBulletCollisionCommon.h"
#include "btBulletDynamicsCommon.h"

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

goRigidBodyHandle goCreateRigidBody(void* user_data, float mass, goCollisionShapeHandle cshape) {
    btCollisionShape* shape = reinterpret_cast<btCollisionShape*>(cshape);
    btAssert(shape);

    btTransform trans;
    trans.setIdentity();

    btVector3 localInertia(0, 0, 0);
    if (mass) {
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

goCollisionShapeHandle goNewSphereShape(goReal radius) { return (goCollisionShapeHandle) new btSphereShape(radius); }

goCollisionShapeHandle goNewBoxShape(goVector3 size) {
    return (goCollisionShapeHandle) new btBoxShape(vec3FromGo(size));
}

goCollisionShapeHandle goNewCylinderShape(goReal radius, goReal height) {
    return (goCollisionShapeHandle) new btCylinderShape(btVector3(radius, height, radius));
}

goCollisionShapeHandle goNewCapsuleShape(goReal radius, goReal height) {
    const int numSpheres = 2;
    btVector3 positions[numSpheres] = {btVector3(0, height, 0), btVector3(0, -height, 0)};
    btScalar radi[numSpheres] = {radius, radius};
    return (goCollisionShapeHandle) new btMultiSphereShape(positions, radi, numSpheres);
}

goCollisionShapeHandle goNewController(goVector3 position, goReal height, goReal stepHeight) {
    auto ghostObject = new btPairCachingGhostObject();
    btTransform transf;
    transf.setIdentity();
    ghostObject->setWorldTransform(transf);

    auto capsule = new btCapsuleShape(0.8, height);
    auto character = new btKinematicCharacterController(ghostObject, capsule, stepHeight);
    return (goCollisionShapeHandle)character;
}

/* Concave static triangle meshes */
goMeshInterfaceHandle goNewMeshInterface() { return 0; }

goCollisionShapeHandle goNewCompoundShape() { return (goCollisionShapeHandle) new btCompoundShape(); }

void goAddChildShape(goCollisionShapeHandle compoundShapeHandle, goCollisionShapeHandle childShapeHandle,
                     goVector3 childPos, goQuaternion childOrn) {
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
// goVector3 v0,goVector3 v1,goVector3 v2); 	extern  goCollisionShapeHandle
// goNewStaticTriangleMeshShape(goMeshInterfaceHandle);

void goAddVertex(goCollisionShapeHandle cshape, goReal x, goReal y, goReal z) {
    btCollisionShape* colShape = reinterpret_cast<btCollisionShape*>(cshape);
    btAssert(colShape->getShapeType() == CONVEX_HULL_SHAPE_PROXYTYPE);

    btConvexHullShape* convexHullShape = reinterpret_cast<btConvexHullShape*>(cshape);
    convexHullShape->addPoint(btVector3(x, y, z));
}

void goDeleteShape(goCollisionShapeHandle cshape) {
    btCollisionShape* shape = reinterpret_cast<btCollisionShape*>(cshape);
    btAssert(shape);
    delete shape;
}

void goSetScaling(goCollisionShapeHandle cshape, goVector3 cscaling) {
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
}

void goGetRotation(goRigidBodyHandle body_handle, goVector3 rotation) {
    btRigidBody* body = reinterpret_cast<btRigidBody*>(body_handle);
    btAssert(body);

    btTransform transf = body->getWorldTransform();
    transf.getRotation().getEulerZYX(rotation[0], rotation[1], rotation[2]);
}

void goSetRotation(goRigidBodyHandle object, const goVector3 rotation) {
    btRigidBody* body = reinterpret_cast<btRigidBody*>(object);
    btAssert(body);

    btQuaternion orient;
    orient.setEuler(rotation[0], rotation[1], rotation[2]);
    btTransform transf = body->getWorldTransform();
    transf.setRotation(orient);
    body->setWorldTransform(transf);
}
