module github.com/johanhenriksson/goworld

go 1.18

require (
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20211213063430-748e38ca8aec
	github.com/go-gl/mathgl v1.0.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/kjk/flex v0.0.0-20171203210503-ed34d6b6a425
	github.com/ojrac/opensimplex-go v1.0.1
	github.com/qmuntal/gltf v0.21.0
	github.com/vulkan-go/vulkan v0.0.0-20210402152248-956e3850d8f9
	github.com/x448/float16 v0.8.4
	golang.org/x/exp v0.0.0-20220323121947-b445f275a754
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/vulkan-go/vulkan => github.com/johanhenriksson/vulkan v0.0.0-20220209212039-bb0a9288948f
