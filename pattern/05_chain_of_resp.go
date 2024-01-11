package pattern

import "fmt"

type Data struct {
	signal             string
	additionalData     string
	rearrangedAsFrames bool
	macAddressAdded    bool
	ipAddressAdded     bool
	transported        bool
	presented          bool
	shownToUser        bool
}

func NewData(signal string) *Data {
	return &Data{
		signal: signal,
	}
}

// Handler interface
type OSILevel interface {
	execute(d *Data)
	setNext(OSILevel)
}

// Concrete handler
type PhysicalLevel struct {
	next OSILevel
}

// Concrete handler method
func (p *PhysicalLevel) execute(d *Data) {
	if d.rearrangedAsFrames {
		fmt.Println("Physical layer already passed")
		p.next.execute(d) // if already processed - immediately send data to the next layer
		return
	}
	fmt.Println("Signal is being processed by the first layer")
	d.additionalData += "Signal framed" // if not processed yet - handle data and send it to the next layer
	d.rearrangedAsFrames = true
	p.next.execute(d)
}

// Setting next handler through method
func (p *PhysicalLevel) setNext(next OSILevel) {
	p.next = next
}

type ChanLevel struct {
	next OSILevel
}

func (c *ChanLevel) GenerateMacAddress(d *Data) {
	d.additionalData += "MAC Address = ..." // if not processed yet - handle data and send it to the next layer
	d.macAddressAdded = true
}

func (c *ChanLevel) execute(d *Data) {
	if d.macAddressAdded {
		fmt.Println("MAC address is already assigned")
		c.next.execute(d) // if already processed - immediately send data to the next layer
		return
	}
	fmt.Println("Adding MAC address")
	c.GenerateMacAddress(d)
	c.next.execute(d)
}

func (c *ChanLevel) setNext(next OSILevel) {
	c.next = next
}

type NetworkLayer struct {
	next      OSILevel
	ipAddress string
}

func (n *NetworkLayer) execute(d *Data) {
	if d.ipAddressAdded {
		fmt.Println("IP Address is already assigned")
		n.next.execute(d) // if already processed - immediately send data to the next layer
		return
	}
	fmt.Println("Adding IP address")
	d.additionalData += fmt.Sprintf("IP address: %s", n.ipAddress)
	d.ipAddressAdded = true
	n.next.execute(d)
}

func (n *NetworkLayer) setNext(next OSILevel) {
	n.next = next
}

type TCPLayer struct {
	next    OSILevel
	tcdData string
}

func (t *TCPLayer) execute(d *Data) {
	if d.transported {
		fmt.Println("Data is already processed with TCP layer")
		t.next.execute(d)
		return
	}
	fmt.Println("Delivering segments")
	d.additionalData += fmt.Sprintf("Data was segmented & delivered")
	d.transported = true
	t.next.execute(d)
}

func (t *TCPLayer) setNext(next OSILevel) {
	t.next = next
}

//... Main idea should be clear at this point

func testChain() {
	// Data to process
	d := NewData("100001001")

	// Handlers  block
	pLayer := &PhysicalLevel{}
	cLayer := &ChanLevel{}
	tLayer := TCPLayer{}
	nLayer := NetworkLayer{ipAddress: "155.83.111.55"}
	// setting order of the chain
	pLayer.setNext(cLayer)
	cLayer.setNext(&nLayer)
	nLayer.setNext(&tLayer)
	// etc

	// sending data to first handler, first handler will take care of it and send to the next handler in chain
	pLayer.execute(d)

}
