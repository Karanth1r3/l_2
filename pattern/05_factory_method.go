package pattern

import "fmt"

// Parameters for factory method
const (
	Knight = iota
	Archer
	Sorcerer
)

// Factory product interface
type Enemy interface {
	attack()
}

// Concrete product
type KnightEnemy struct {
	attackPower int
}

func NewKnightEnemy() *KnightEnemy {
	return &KnightEnemy{attackPower: 50}
}

// Implementing product interface
func (k KnightEnemy) attack() {
	fmt.Printf("Attacked by knight, received %d damage\n", k.attackPower)
}

// Product 2
type ArcherEnemy struct {
	attackPower, attackRange int
}

func NewArcherEnemy() *ArcherEnemy {
	return &ArcherEnemy{attackPower: 20, attackRange: 50}
}

func (a ArcherEnemy) attack() {
	fmt.Printf("Attacked by archer, received %d damage from %d range\n", a.attackPower, a.attackRange)
}

// Product 3
type SorcererEnemy struct {
	attackPower, attackRange, debuffDuration int
}

func NewSorcererEnemy() *SorcererEnemy {
	return &SorcererEnemy{attackPower: 10, attackRange: 75, debuffDuration: 5}
}

func (s SorcererEnemy) attack() {
	fmt.Printf("Attacked by sorcerer, received %d damage from %d range. Receiving dmg from debuff for next %d seconds\n", s.attackPower, s.attackRange, s.debuffDuration)
}

// returning interface to avoid too deep connection with concrete implementations
func CreateEnemy(enemyType int) (Enemy, error) {
	switch enemyType {
	case Knight:
		return NewKnightEnemy(), nil
	case Archer:
		return NewArcherEnemy(), nil
	case Sorcerer:
		return NewSorcererEnemy(), nil
	}
	return nil, fmt.Errorf("Wrong enemy type passed")
}

func testFactory() {
	knight, _ := CreateEnemy(Knight) // logic of creating new objects is stored in factory method now now
	sorcerer, _ := CreateEnemy(Sorcerer)
	archer, _ := CreateEnemy(Archer)

	knight.attack()
	sorcerer.attack()
	archer.attack()
}
