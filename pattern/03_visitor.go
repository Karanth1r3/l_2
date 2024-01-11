package pattern

import "fmt"

// New functionality will be added through these funcs
type visitor interface {
	visitSeniorStaff(*SeniorStaff)
	visitMiddleStaff(*MiddleStaff)
	visitJuniorStaff(*JuniorStuff)
}

// Concrete visitor, struct with methods whcih add new functionality "outside" of original structures
type WorkLoadCalculator struct {
	workLoadIndex int
}

func (swlc *WorkLoadCalculator) visitJuniorStaff(j *JuniorStuff) {
	loadIndex := j.productivityIndex * 999
	swlc.workLoadIndex = loadIndex
	fmt.Printf("Junior staff workload is calculated as %d", swlc.workLoadIndex)
}

func (swlc *WorkLoadCalculator) visitMiddleStaff(m *MiddleStaff) {
	loadIndex := m.productivityIndex * 2500
	swlc.workLoadIndex = loadIndex
	fmt.Printf("Middle staff workload is calculated as %d", swlc.workLoadIndex)
}

func (swlc *WorkLoadCalculator) visitSeniorStaff(s *SeniorStaff) {
	loadIndex := s.productivityIndex * 1
	swlc.workLoadIndex = loadIndex
	fmt.Printf("Senior staff workload is calculated as %d", swlc.workLoadIndex)
}

// Interface for types in which new functionality should be added
type Staff interface {
	getType() string
	accept(visitor)
}

// Concrete type to add functionality
type SeniorStaff struct {
	productivityIndex int
}

// accept visitor for choosing appropriate function frov visitor based on accepter type
func (s *SeniorStaff) accept(v visitor) {
	v.visitSeniorStaff(s)
}

func (s *SeniorStaff) getType() string {
	return "SeniorStaff"
}

// Another one
type MiddleStaff struct {
	productivityIndex int
}

func (m *MiddleStaff) accept(v visitor) {
	v.visitMiddleStaff(m)
}

func (m *MiddleStaff) getType() string {
	return "MiddleStaff"
}

// And another one
type JuniorStuff struct {
	productivityIndex int
}

func (j *JuniorStuff) accept(v visitor) {
	v.visitJuniorStaff(j)
}

func (j *JuniorStuff) getType() string {
	return "JuniorStaff"
}

func testVisitor() {
	//objects with new functionality
	jun := &JuniorStuff{productivityIndex: 1}
	mid := &MiddleStaff{productivityIndex: 3}
	sen := &SeniorStaff{productivityIndex: 5}
	// new functionality struct
	workLoadCalc := &WorkLoadCalculator{}
	// calling new functionality
	jun.accept(workLoadCalc)
	mid.accept(workLoadCalc)
	sen.accept(workLoadCalc)
}
