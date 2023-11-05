package physics

/*
#cgo CXXFLAGS: -std=c++11 -I/usr/local/include/bullet
#cgo CFLAGS: -I/usr/local/include/bullet
#cgo LDFLAGS: -lstdc++ -L/usr/local/lib -lBulletDynamics -lBulletCollision -lLinearMath -lBullet3Common
#include "bullet.h"
*/
import "C"

import (
	"reflect"
	"unsafe"

	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/vertex"
)

//
// utilities
//

func init() {
	// sanity check

	vsize := unsafe.Sizeof(vec3.T{})
	if vsize != 12 {
		panic("expected vec3 to be 12 bytes")
	}

	qsize := unsafe.Sizeof(quat.T{})
	if qsize != 16 {
		panic("expected quaternion to be 16 bytes")
	}
}

func vec3ptr(v *vec3.T) *C.goVector3 {
	return (*C.goVector3)(unsafe.Pointer(v))
}

func quatPtr(q *quat.T) *C.goQuaternion {
	return (*C.goQuaternion)(unsafe.Pointer(q))
}

//
// dynamics world
//

type worldHandle C.goDynamicsWorldHandle

func world_new() worldHandle {
	handle := C.goCreateDynamicsWorld()
	return worldHandle(handle)
}

func world_delete(world *worldHandle) {
	C.goDeleteDynamicsWorld(*world)
	*world = nil
}

func world_add_rigidbody(world worldHandle, body rigidbodyHandle, group, mask Mask) {
	C.goAddRigidBody(world, body, C.int(group), C.int(mask))
}

func world_remove_rigidbody(world worldHandle, body rigidbodyHandle) {
	C.goRemoveRigidBody(world, body)
}

func world_add_character(world worldHandle, character characterHandle) {
	C.goAddCharacter(world, character)
}

func world_remove_character(world worldHandle, character characterHandle) {
	C.goRemoveCharacter(world, character)
}

func world_debug_draw(world worldHandle) {
	C.goDebugDraw(world)
}

func world_step_simulation(world worldHandle, timestep, fixedTimestep float32, maxSteps int) {
	C.goStepSimulation(world, (C.goReal)(timestep), (C.goReal)(fixedTimestep), (C.int)(maxSteps))
}

func world_gravity_set(world worldHandle, gravity vec3.T) {
	C.goSetGravity(world, vec3ptr(&gravity))
}

func world_debug_enable(w *World) {
	C.goEnableDebug(w.handle)
}

func world_debug_disable(w *World) {
	C.goDisableDebug(w.handle)
}

type raycastResult struct {
	shape  unsafe.Pointer
	point  vec3.T
	normal vec3.T
}

func world_raycast(world worldHandle, from, to vec3.T, mask Mask) (result raycastResult, hit bool) {
	hits := C.goRayCast(world, vec3ptr(&from), vec3ptr(&to), C.int(mask), (*C.goRayCastResult)(unsafe.Pointer(&result)))
	if hits > 0 {
		hit = true
	}
	return
}

//
// character functions
//

type characterHandle C.goCharacterHandle

type characterState struct {
	position vec3.T
	rotation quat.T
	grounded bool
}

func character_new(shape shapeHandle, stepHeight float32) characterHandle {
	handle := C.goCreateCharacter(shape, C.float(stepHeight))
	return characterHandle(handle)
}

func character_delete(char *characterHandle) {
	C.goDeleteCharacter(*char)
	*char = nil
}

func character_state_pull(char characterHandle) characterState {
	state := characterState{}
	C.goCharacterGetState(char, (*C.goCharacterState)(unsafe.Pointer(&state)))
	return state
}

func character_state_push(char characterHandle, position vec3.T, rotation quat.T) {
	state := characterState{
		position: position,
		rotation: rotation,
	}
	C.goCharacterSetState(char, (*C.goCharacterState)(unsafe.Pointer(&state)))
}

func character_move(char characterHandle, dir vec3.T) {
	C.goCharacterMove(char, vec3ptr(&dir))
}

func character_jump(char characterHandle) {
	C.goCharacterJump(char)
}

//
// rigidbody
//

type rigidbodyHandle C.goRigidBodyHandle

type rigidbodyState struct {
	position vec3.T
	rotation quat.T
	mass     float32
}

func rigidbody_new(ptr unsafe.Pointer, mass float32, shape shapeHandle) rigidbodyHandle {
	handle := C.goCreateRigidBody((*C.char)(ptr), C.goReal(mass), shape)
	return rigidbodyHandle(handle)
}

func rigidbody_delete(body *rigidbodyHandle) {
	C.goDeleteRigidBody(*body)
	*body = nil
}

func rigidbody_shape_set(body rigidbodyHandle, shape shapeHandle) {
	C.goRigidBodySetShape(body, shape)
}

func rigidbody_state_pull(body rigidbodyHandle) rigidbodyState {
	state := rigidbodyState{}
	C.goRigidBodyGetState(body, (*C.goRigidBodyState)(unsafe.Pointer(&state)))
	return state
}

func rigidbody_state_push(body rigidbodyHandle, position vec3.T, rotation quat.T) {
	state := rigidbodyState{
		position: position,
		rotation: rotation,
	}
	C.goRigidBodySetState(body, (*C.goRigidBodyState)(unsafe.Pointer(&state)))
}

//
// shape
//

type shapeHandle C.goShapeHandle

func shape_new_box(ptr unsafe.Pointer, extents vec3.T) shapeHandle {
	handle := C.goNewBoxShape((*C.char)(ptr), vec3ptr(&extents))
	return shapeHandle(handle)
}

func shape_new_capsule(ptr unsafe.Pointer, radius, height float32) shapeHandle {
	handle := C.goNewCapsuleShape((*C.char)(ptr), C.float(radius), C.float(height))
	return shapeHandle(handle)
}

func shape_new_sphere(ptr unsafe.Pointer, radius float32) shapeHandle {
	handle := C.goNewSphereShape((*C.char)(ptr), C.float(radius))
	return shapeHandle(handle)
}

func shape_new_triangle_mesh(ptr unsafe.Pointer, mesh meshHandle) shapeHandle {
	handle := C.goNewTriangleMeshShape((*C.char)(ptr), mesh)
	return shapeHandle(handle)
}

func shape_new_compound(ptr unsafe.Pointer) shapeHandle {
	handle := C.goNewCompoundShape((*C.char)(ptr))
	return shapeHandle(handle)
}

func shape_scaling_set(shape shapeHandle, scale vec3.T) {
	C.goSetScaling(shape, vec3ptr(&scale))
}

func compound_add_child(shape, child shapeHandle, position vec3.T, rotation quat.T) {
	C.goAddChildShape(shape, child, vec3ptr(&position), quatPtr(&rotation))
}

func compound_update_child(shape shapeHandle, index int, position vec3.T, rotation quat.T) {
	C.goUpdateChildShape(shape, C.int(index), vec3ptr(&position), quatPtr(&rotation))
}

func compound_remove_child(shape, child shapeHandle) {
	C.goRemoveChildShape(shape, child)
}

func shape_delete(shape *shapeHandle) {
	C.goDeleteShape(*shape)
	*shape = nil
}

//
// mesh
//

type meshHandle C.goTriangleMeshHandle

func mesh_new(mesh vertex.Mesh) meshHandle {
	vertexArray := reflect.ValueOf(mesh.VertexData())
	if vertexArray.Kind() != reflect.Slice {
		panic("vertex data is not a slice")
	}
	if vertexArray.Len() < 1 {
		panic("vertex data is empty")
	}
	vertexPtr := vertexArray.Index(0).UnsafeAddr()
	vertexCount := vertexArray.Len()

	indexArray := reflect.ValueOf(mesh.IndexData())
	if indexArray.Kind() != reflect.Slice {
		panic("index data is not a slice")
	}
	if indexArray.Len() < 1 {
		panic("index data is empty")
	}
	indexPtr := indexArray.Index(0).UnsafeAddr()
	indexCount := indexArray.Len()

	handle := C.goNewTriangleMesh(
		unsafe.Pointer(vertexPtr), C.int(vertexCount), C.int(mesh.VertexSize()),
		unsafe.Pointer(indexPtr), C.int(indexCount), C.int(mesh.IndexSize()))
	return meshHandle(handle)
}

func mesh_delete(mesh *meshHandle) {
	C.goDeleteTriangleMesh(*mesh)
	*mesh = nil
}
