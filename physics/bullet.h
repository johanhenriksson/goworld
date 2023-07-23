#ifndef BULLET_C_API_H
#define BULLET_C_API_H

#include <stdbool.h>

#define GO_DECLARE_HANDLE(name) \
    typedef struct name##__ {   \
        int unused;             \
    }* name

typedef float goReal;

typedef struct {
    goReal x;
    goReal y;
    goReal z;
} goVector3;

typedef struct {
    goReal w;
    goReal x;
    goReal y;
    goReal z;
} goQuaternion;

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

extern void goSetGravity(goDynamicsWorldHandle world, goVector3* gravity);

extern void goDeleteDynamicsWorld(goDynamicsWorldHandle world);

extern void goStepSimulation(goDynamicsWorldHandle, goReal timeStep);

extern void goAddRigidBody(goDynamicsWorldHandle world, goRigidBodyHandle object);

extern void goRemoveRigidBody(goDynamicsWorldHandle world, goRigidBodyHandle object);

/* Rigid Body  */

typedef struct {
    goVector3 position;
    goQuaternion rotation;
    goReal mass;
} goRigidBodyState;

extern goRigidBodyHandle goCreateRigidBody(char* user_data, float mass, goShapeHandle cshape);

extern void goDeleteRigidBody(goRigidBodyHandle body);

extern void goRigidBodyGetState(goRigidBodyHandle body, goRigidBodyState* state);
extern void goRigidBodySetState(goRigidBodyHandle body, goRigidBodyState* state);
extern void goRigidBodySetShape(goRigidBodyHandle objectPtr, goShapeHandle shapePtr);

/* Collision Shape definition */

extern goShapeHandle goNewSphereShape(char* user, goReal radius);

extern goShapeHandle goNewBoxShape(char* user, goVector3* size);

extern goShapeHandle goNewCapsuleShape(char* user, goReal radius, goReal height);

extern goShapeHandle goNewCompoundShape(char* user);

extern void goAddChildShape(goShapeHandle compoundShape, goShapeHandle childShape, goVector3* pos, goQuaternion* rot);

extern void goUpdateChildShape(goShapeHandle compoundShape, int index, goVector3* pos, goQuaternion* rot);

extern void goRemoveChildShape(goShapeHandle compoundShape, goShapeHandle childShape);

extern void goDeleteShape(goShapeHandle shape);

/* Convex Meshes */
extern goShapeHandle goNewConvexHullShape(void);
extern void goAddVertex(goShapeHandle convexHull, goReal x, goReal y, goReal z);

/* Concave static triangle meshes */
extern goTriangleMeshHandle goNewTriangleMesh(void* vertex_ptr, int vertex_count, int vertex_stride, void* index_ptr,
                                              int index_count, int index_stride);

extern void goDeleteTriangleMesh(goTriangleMeshHandle handle);

extern goShapeHandle goNewTriangleMeshShape(char* user, goTriangleMeshHandle);

extern void goSetScaling(goShapeHandle shape, goVector3* scaling);

// raycast

typedef struct goRayCastResult {
    void* shape;
    goVector3 point;
    goVector3 normal;
} goRayCastResult;

extern int goRayCast(goDynamicsWorldHandle world, goVector3* rayStart, goVector3* rayEnd, goRayCastResult* res);

// debugging

extern void goEnableDebug(goDynamicsWorldHandle world);
extern void goDisableDebug(goDynamicsWorldHandle world);

extern void goDebugDraw(goDynamicsWorldHandle world);

// character controller

typedef struct {
    goVector3 position;
    goQuaternion rotation;
    bool grounded;
} goCharacterState;

extern goCharacterHandle goCreateCharacter(goShapeHandle shapeHandle, float stepHeight);
extern void goDeleteCharacter(goCharacterHandle handle);
extern void goCharacterMove(goCharacterHandle handle, goVector3* direction);
extern void goCharacterJump(goCharacterHandle handle);
extern void goAddCharacter(goDynamicsWorldHandle world, goCharacterHandle handle);
extern void goRemoveCharacter(goDynamicsWorldHandle world, goCharacterHandle handle);

extern void goCharacterSetState(goCharacterHandle handle, goCharacterState* state);
extern void goCharacterGetState(goCharacterHandle handle, goCharacterState* state);

#ifdef __cplusplus
}
#endif

#endif  // BULLET_C_API_H
