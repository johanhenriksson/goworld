
#include <functional>

#include "LinearMath/btIDebugDraw.h"
#include "LinearMath/btVector3.h"
#include "bullet.h"

typedef void (*DrawLineCallback)(const btVector3&, const btVector3&, const btVector3&);

class GoDebugDrawer : public btIDebugDraw {
    int m_debugMode;

    std::function<void(const btVector3&, const btVector3&, const btVector3&)> drawline;

   public:
    GoDebugDrawer(std::function<void(const btVector3&, const btVector3&, const btVector3&)> drawline_cb) {
        // callbacks
        this->drawline = drawline_cb;

        // default drawing mode
        this->m_debugMode = DBG_DrawAabb;
    }

    virtual void drawLine(const btVector3& from, const btVector3& to, const btVector3& color) {
        // Here is where you would actually draw the line.
        // This could be done using OpenGL, DirectX, or any other graphics API.
        // The parameters "from" and "to" specify the endpoints of the line, and "color" specifies its color.
        this->drawline(from, to, color);
    }

    virtual void reportErrorWarning(const char* warningString) { printf("PHYS WARNING: %s\n", warningString); }

    virtual void setDebugMode(int debugMode) { m_debugMode = debugMode; }

    virtual int getDebugMode() const { return m_debugMode; }

    virtual void drawContactPoint(const btVector3& PointOnB, const btVector3& normalOnB, btScalar distance,
                                  int lifeTime, const btVector3& color) {}

    virtual void draw3dText(const btVector3& location, const char* textString) {}
};
