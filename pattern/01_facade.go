package pattern

import "fmt"

// renderer modules, for basic operations with complex system it's fine to create
// a facade class
type RenderPipelineFacade struct {
	vertexProcessing   *VertexProcessing
	clipping           *Clipping
	rasterizer         *Rasterizer
	fragmentProcessing *FragmentProcessing
	postProcessing     *PostProcessing
}

const (
	Mobile int = iota
	Console
	Desktop
)

func newRPFacade(nearClipPlane, farClipPlane float64, postProcessingSamplingTier int) *RenderPipelineFacade {
	fmt.Println("Rendering Pipeline Facade initializing")
	renderPipelineFacade := &RenderPipelineFacade{
		vertexProcessing:   &VertexProcessing{},
		clipping:           newClipping(nearClipPlane, farClipPlane),
		rasterizer:         &Rasterizer{},
		fragmentProcessing: &FragmentProcessing{},
		postProcessing:     newPostProcessingModule(Desktop),
	}
	return renderPipelineFacade
}

// Wrapping complex system modules functions through facade functions
func (rpf *RenderPipelineFacade) ModifyRasterizer(v VertexData) {
	rpf.rasterizer.customize(v)
}

func (rpf *RenderPipelineFacade) CallRasterize(v VertexData) { // technically it's ridicilous but who cares
	rpf.rasterizer.rasterize(v)
}

func (rpf *RenderPipelineFacade) SetPostProcessingState(state bool) {
	rpf.postProcessing.setPostProcessingState(state)
}

func (rpf *RenderPipelineFacade) UpdateVertexData() (v *VertexData) {
	vData := rpf.vertexProcessing.CreateVertexData()
	return vData
}

// modules of the complex system, required interactions with them are placed in
// a single facade "class" (one place)
type VertexProcessing struct {
	vertices []float64
}

type VertexData struct {
	vData []float64
}

func (v *VertexProcessing) CreateVertexData() (vd *VertexData) {
	vd = &VertexData{
		vData: []float64{0.1, 0.5, 1},
	}
	fmt.Println("Vertex Processing stage is done")
	return vd
}

type Clipping struct {
	near, far float64
}

// Constructor
func newClipping(near, far float64) *Clipping {
	clipping := &Clipping{
		near: near,
		far:  far,
	}
	return clipping
}

type Rasterizer struct {
	rData []float64
}

func (r *Rasterizer) rasterize(v VertexData) {
	fmt.Println("Rasterization stage is done")
}

func (r *Rasterizer) customize(v VertexData) {
	for _, elem := range v.vData {
		elem += 0.5
	}
	fmt.Println("Raster module was modified")
}

type FragmentProcessing struct {
	fData []int
}

type PostProcessing struct {
	quality int
	enabled bool
}

func newPostProcessingModule(q int) *PostProcessing {
	postProcessing := &PostProcessing{
		quality: q,
	}
	fmt.Println("Post Processing initialized with parameters according to hardware specs")
	return postProcessing
}

func (p *PostProcessing) setPostProcessingState(state bool) {
	p.enabled = state
}

func main() {
	// Initializing facade with constructor
	renderingPipelineFacade := newRPFacade(0.1, 1000, Desktop)
	// Using complex system modules through facade tools, "API"
	vertexData := renderingPipelineFacade.UpdateVertexData()
	renderingPipelineFacade.ModifyRasterizer(*vertexData) // doesn't make sense technically but who cares
	renderingPipelineFacade.CallRasterize(*vertexData)    // i guess i should've been more aware of the rendering pipeline guts at this point to make up better functions examples but i am not
	//

}

// Facade pattern is useful to wrap a complex system functional in one place.
// Client can learn how to use the facade with its (probably limited but not necessary) functionality instead of learning the whole complex system & its functions.

// It is required to know how the system shall (and should ideally) be used to properly implement this pattern
// "API" of the facade may not fully represent system capabilities but should be enough to use it within the required task context

// There is risk of making the facade a GodComponent (an anti-pattern) but if used mindfully - it's a fine way to isolate interaction with a complex system within a single "class"
// If a single facade is becoming too large => it's possible to create additional facades

// example of real-conditions usage - implement limited api for game engine renderer modules for artists/technical artists to easily create custom graphics effects
