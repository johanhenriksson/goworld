#ifndef BULLET_C_API_H
#define BULLET_C_API_H

#define GO_DECLARE_HANDLE(name) \
    typedef struct name##__ {   \
        int unused;             \
    }* name

typedef float goReal;

typedef goReal goVector3[3];
typedef goReal goQuaternion[4];

#ifdef __cplusplus
extern "C" {
#endif

// 	Dynamics world, belonging to some physics SDK
GO_DECLARE_HANDLE(goDynamicsWorldHandle);

// Rigid Body that can be part of a Dynamics World
GO_DECLARE_HANDLE(goRigidBodyHandle);

// 	Collision Shape/Geometry, property of a Rigid Body
GO_DECLARE_HANDLE(goShapeHandle);

// Constraint for Rigid Bodies */
GO_DECLARE_HANDLE(goConstraintHandle);

// Triangle Mesh interface
GO_DECLARE_HANDLE(goMeshInterfaceHandle);
GO_DECLARE_HANDLE(goTriangleMeshHandle);

// Broadphase Scene/Proxy Handles
GO_DECLARE_HANDLE(goCollisionBroadphaseHandle);
GO_DECLARE_HANDLE(goBroadphaseProxyHandle);
GO_DECLARE_HANDLE(goCollisionWorldHandle);

GO_DECLARE_HANDLE(goCharacterHandle);

typedef void (*btBroadphaseCallback)(void* clientData, void* object1, void* object2);

extern goCollisionBroadphaseHandle goCreateSapBroadphase(btBroadphaseCallback beginCallback,
                                                         btBroadphaseCallback endCallback);

extern void goDestroyBroadphase(goCollisionBroadphaseHandle bp);

extern goBroadphaseProxyHandle goCreateProxy(goCollisionBroadphaseHandle bp, void* clientData, goReal minX, goReal minY,
                                             goReal minZ, goReal maxX, goReal maxY, goReal maxZ);

extern void goDestroyProxy(goCollisionBroadphaseHandle bp, goBroadphaseProxyHandle proxyHandle);

extern void goSetBoundingBox(goBroadphaseProxyHandle proxyHandle, goReal minX, goReal minY, goReal minZ, goReal maxX,
                             goReal maxY, goReal maxZ);

/* todo: add pair cache support with queries like add/remove/find pair */

/* Dynamics World */

extern goDynamicsWorldHandle goCreateDynamicsWorld();

extern void goSetGravity(goDynamicsWorldHandle world, goVector3 gravity);

extern void goDeleteDynamicsWorld(goDynamicsWorldHandle world);

extern void goStepSimulation(goDynamicsWorldHandle, goReal timeStep);

extern void goAddRigidBody(goDynamicsWorldHandle world, goRigidBodyHandle object);

extern void goRemoveRigidBody(goDynamicsWorldHandle world, goRigidBodyHandle object);

/* Rigid Body  */

extern goRigidBodyHandle goCreateRigidBody(void* user_data, float mass, goShapeHandle cshape);

extern void goDeleteRigidBody(goRigidBodyHandle body);

/* Collision Shape definition */

extern goShapeHandle goNewSphereShape(goReal radius);

extern goShapeHandle goNewBoxShape(goVector3 size);

extern goShapeHandle goNewCapsuleShape(goReal radius, goReal height);

extern goShapeHandle goNewCompoundShape(void);

extern void goAddChildShape(goShapeHandle compoundShape, goShapeHandle childShape, goVector3 childPos,
                            goQuaternion childOrn);

extern void goDeleteShape(goShapeHandle shape);

/* Convex Meshes */
extern goShapeHandle goNewConvexHullShape(void);
extern void goAddVertex(goShapeHandle convexHull, goReal x, goReal y, goReal z);

/* Concave static triangle meshes */
extern goTriangleMeshHandle goNewTriangleMesh(void* vertex_ptr, int vertex_count, int vertex_stride, void* index_ptr,
                                              int index_count, int index_stride);

extern goShapeHandle goNewStaticTriangleMeshShape(goTriangleMeshHandle);

extern void goSetScaling(goShapeHandle shape, goVector3 scaling);

/* get world transform */
extern void goGetPosition(goRigidBodyHandle object, goVector3 position);
extern void goGetRotation(goRigidBodyHandle object, goVector3 rotation);

/* set world transform (position/orientation) */
extern void goSetPosition(goRigidBodyHandle object, const goVector3 position);
extern void goSetRotation(goRigidBodyHandle object, const goVector3 rotation);

// raycast

typedef struct goRayCastResult {
    goRigidBodyHandle m_body;
    goShapeHandle m_shape;
    goVector3 m_positionWorld;
    goVector3 m_normalWorld;
} goRayCastResult;

extern int goRayCast(goDynamicsWorldHandle world, const goVector3 rayStart, const goVector3 rayEnd,
                     goRayCastResult res);

// debugging

extern void goEnableDebug(goDynamicsWorldHandle world);
extern void goDisableDebug(goDynamicsWorldHandle world);

extern void goDebugDraw(goDynamicsWorldHandle world);

// character controller

extern goCharacterHandle goCreateCharacter(goShapeHandle shapeHandle, float height, float radius, float stepHeight);
extern void goDeleteCharacter(goCharacterHandle handle);
extern void goCharacterWalkDirection(goCharacterHandle handle, goVector3 direction);
extern void goCharacterJump(goCharacterHandle handle);
extern void goCharacterWarp(goCharacterHandle handle, goVector3 direction);
extern void goCharacterUpdate(goCharacterHandle handle, goDynamicsWorldHandle world, float dt);
extern void goAddCharacter(goDynamicsWorldHandle world, goCharacterHandle handle);
extern void goRemoveCharacter(goDynamicsWorldHandle world, goCharacterHandle handle);

#ifdef __cplusplus
}
#endif

#endif  // BULLET_C_API_H
