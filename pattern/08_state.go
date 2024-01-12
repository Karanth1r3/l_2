package pattern

import "fmt"

type Character struct {
	inCutscene CharState
	playable   CharState
	autoMoving CharState

	currentState CharState

	movingControllerActive bool
	autoMovingActive       bool
	cameraControllerActive bool

	Position
}

type Position struct {
	x, y int
}

func NewCharacter(initPos Position) *Character {
	ch := &Character{
		movingControllerActive: true,
		autoMovingActive:       false,
		cameraControllerActive: true,

		Position: initPos,
	}

	inCutScene := &CutsceneState{
		character: ch,
	}

	playable := &PlayableState{
		character: ch,
	}

	auto := &AutoMoveState{
		character: ch,
	}

	ch.inCutscene = inCutScene
	ch.playable = playable
	ch.autoMoving = auto

	ch.setState(playable)
	return ch
}

func (c *Character) setState(s CharState) {
	c.currentState = s
}

func (c *Character) enableControl() {
	c.movingControllerActive = true
	c.autoMovingActive = false
	c.cameraControllerActive = true
}

func (c *Character) disableControl() {
	c.movingControllerActive = false
	c.autoMovingActive = false
	c.cameraControllerActive = false
}

func (c *Character) setAutoMove() {
	c.movingControllerActive = false
	c.autoMovingActive = true
	c.cameraControllerActive = true
}

func (c *Character) move(p Position) {
	c.Position = p
}

func (c *Character) startCutscene() {
	c.currentState.startCutscene()
}

func (c *Character) stopCutscene() {
	c.currentState.stopCutscene()
}

func (c *Character) autoMove() {
	fmt.Println("Auto moving to target position")
}

// State interface
type CharState interface {
	move(pos Position)
	autoMove()
	startCutscene()
	stopCutscene()
}

// Concrete State 1
type CutsceneState struct {
	character *Character
}

// Behaviour of characted while in state 1. In cutscene characted will not be able to move, for example
func (c *CutsceneState) move(pos Position) {
	fmt.Println("Cannot move while in cutscene")
}

func (c *CutsceneState) autoMove() {
	fmt.Println("Auto move module was deactivated due to cutscene commence")
}

func (c *CutsceneState) startCutscene() {
	c.character.disableControl()
	fmt.Println("Control was disabled due to cutscene commence")
}

func (c *CutsceneState) stopCutscene() {
	c.character.setState(c.character.playable)
	c.character.enableControl()
	fmt.Println("Control granted back")
}

// Concrete state 2
type PlayableState struct {
	character *Character
}

func (p *PlayableState) move(pos Position) {
	p.character.enableControl()
	p.character.move(pos)
}

func (p *PlayableState) autoMove() {
	p.character.setState(p.character.autoMoving)
	p.character.setAutoMove()
}

func (p *PlayableState) startCutscene() {
	p.character.setState(p.character.inCutscene)
}

func (p *PlayableState) stopCutscene() {
	fmt.Println("Control is going to be granted back")
}

// Concrete state 3
type AutoMoveState struct {
	character *Character
}

func (a *AutoMoveState) move(pos Position) {
	fmt.Println("Move is not called while auto move is active")
}

func (a *AutoMoveState) autoMove() {
	a.character.autoMove()
}

func (a *AutoMoveState) startCutscene() {
	a.character.setState(a.character.inCutscene)
}

func (a *AutoMoveState) stopCutscene() {
	fmt.Println("Should do nothing here")
}

func testState() {
	hero := NewCharacter(Position{0, 0})

	hero.move(Position{15, 5}) // Operations from moving scene are executing

	hero.startCutscene()      // Changing state to Cutscene
	hero.move(Position{0, 0}) // In cutscene state hero won't move

	hero.stopCutscene() // Going back to initial state

	hero.setState(hero.autoMoving) // States can be set like that too

	hero.stopCutscene() // Does nothing in current state

	hero.setState(hero.playable) //Initial state

	hero.move(Position{0, 0}) // Back to home
}
