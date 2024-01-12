package pattern

import "fmt"

// Strategy interface
type Robbery interface {
	Rob(*Bank)
}

type CriminalWallet struct {
	money int
}

// Concrete strategy 1
type StealthRobbery struct {
	CriminalWallet
	crewStealthMastery int
	effectiveness      int
}

// Functionality which can be switched easily with another concrete strategy implementation in service
func (s *StealthRobbery) Rob(b *Bank) {
	if s.crewStealthMastery > b.awareness {
		b.money -= s.effectiveness
		s.money += s.effectiveness
	} else {
		s.crewStealthMastery -= 10
		s.effectiveness -= 10
		b.defenseSystemIndex += 5
		b.awareness += 20
		fmt.Println("Stealth robbery failed. Bank security systems have been improved")
	}
}

type ArmedRobbery struct {
	CriminalWallet
	crewGunfightMastery int
	effectiveness       int
}

func (a *ArmedRobbery) Rob(b *Bank) {
	if a.crewGunfightMastery > b.defenseSystemIndex {
		b.money -= a.effectiveness
		a.money += a.effectiveness
		fmt.Println("Armed robbery was successfull")
	} else {
		a.crewGunfightMastery -= 10
		a.effectiveness -= 10
		b.defenseSystemIndex += 20
		b.awareness += 5
		fmt.Println("Armed robbery failed. Bank security systems have been improved")
	}
}

// Service
type CriminalCrew struct {
	robStrat Robbery
}

// Service can change behaviour strategy by assigning different ones
func (c *CriminalCrew) setRobberyStrategy(r Robbery) {
	c.robStrat = r
}

func (c *CriminalCrew) RobBank(b *Bank) {
	c.robStrat.Rob(b)
}

// This one actually is not connected to strategy pattern but in this example it is necessary
type Bank struct {
	money, defenseSystemIndex, awareness int
}

func NewBank(money, defense, awareness int) *Bank {
	return &Bank{
		money:              money,
		defenseSystemIndex: defense,
		awareness:          awareness,
	}
}

func testStrategy() {
	bank := NewBank(50000, 15, 25)

	wallet := &CriminalWallet{0}

	armedRob := ArmedRobbery{*wallet, 50, 750}
	stealthRob := StealthRobbery{*wallet, 24, 1000}

	// Trying with the first strategy
	crew := CriminalCrew{&stealthRob}
	crew.RobBank(bank) // Should end as fail
	// Changing strategy
	crew.setRobberyStrategy(&armedRob)
	crew.RobBank(bank) // This one should be successfull

}
