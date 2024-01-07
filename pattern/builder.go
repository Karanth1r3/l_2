package pattern

import "fmt"

const (
	Gaming = iota
	Workstation
)

// Composite structure/ class which includes many different elements/parts
type PC struct {
	cpu cpu
	gpu gpu
	ram ram
}

type gpu struct {
	info   string
	tflops float64
}

type cpu struct {
	info      string
	coreCount int
}

type ram int

func (pc PC) TestPC() {
	fmt.Printf("\nPC is build & working %v, %v and %v, are functioning properly\n", pc.cpu, pc.gpu, pc.ram)
}

// Telescopic/large constructors can be replaced with builder based on this interace
type IBuilder interface {
	setCPU()
	setGPU()
	setRAM()
	getPC() *PC
}

func getPCBuilder(PCType int) IBuilder {

	switch PCType {

	case Gaming:
		return newGamingPCBuilder()
	case Workstation:
		return newWorkstationPCBuilder()
	}
	return nil
}

// First concrete builder
type gamingPCBuilder struct {
	cpu cpu
	gpu gpu
	ram ram
}

// Builder constructor
func newGamingPCBuilder() *gamingPCBuilder {
	return &gamingPCBuilder{}
}

// Components setup methods
func (g *gamingPCBuilder) setCPU() {
	g.cpu = cpu{"Intel Core", 16}
}

func (g *gamingPCBuilder) setGPU() {
	g.gpu = gpu{"AMD Radeon", 70}
}

func (g *gamingPCBuilder) setRAM() {
	g.ram = 32
}

// getting builder result method
func (g *gamingPCBuilder) getPC() *PC {
	return &PC{g.cpu, g.gpu, g.ram}
}

// Second concrete builder with same semantics
type workstationPCBuilder struct {
	cpu cpu
	gpu gpu
	ram ram
}

func newWorkstationPCBuilder() *workstationPCBuilder {
	return &workstationPCBuilder{}
}

func (w *workstationPCBuilder) setCPU() {
	w.cpu = cpu{"AMD Threadripper", 128}
}

func (w *workstationPCBuilder) setGPU() {
	w.gpu = gpu{"Nvidia Tesla", 70}
}

func (w *workstationPCBuilder) setRAM() {
	w.ram = 32
}

func (w *workstationPCBuilder) getPC() *PC {
	return &PC{w.cpu, w.gpu, w.ram}
}

// Client can set any concrete builder as active to the director. then directors current builder can be called to build required object
type director struct {
	builder IBuilder
}

func newDirector(b IBuilder) *director {
	return &director{
		builder: b,
	}
}

func (d *director) setBuilder(b IBuilder) {
	d.builder = b
}

func (d *director) buildPC() PC {
	d.builder.setCPU()
	d.builder.setGPU()
	d.builder.setRAM()
	return *d.builder.getPC()
}

func builderFunc() {
	// initializing builders
	wsBuilder := getPCBuilder(Workstation)
	gamingPCBuilder := getPCBuilder(Gaming)
	// creating director with assigned concrete builder & building object with it
	dir := newDirector(wsBuilder)
	workStation := dir.buildPC()
	workStation.TestPC()
	// assigning different concrete builder & building second objects variation
	dir.setBuilder(gamingPCBuilder)
	gamingPC := dir.buildPC()
	gamingPC.TestPC()
}

//
