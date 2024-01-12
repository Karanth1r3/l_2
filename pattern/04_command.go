package pattern

import "fmt"

// sender, command should be layer between sender and logic elements
type Button struct {
	command Command
}

func (b *Button) press() {
	b.command.execute()
}

// Simplified app struct
type Applictation struct {
	text         string
	clipboard    string
	editors      []Editor
	activeEditor Editor
	history      CommandHistory
}

func NewApplication(e *Editor) *Applictation {
	return &Applictation{
		clipboard:    "",
		activeEditor: *e,
		editors:      []Editor{*e},
		history:      CommandHistory{},
	}
}

// App executes command sent from ui elements
func (app *Applictation) executeCommand(c Command) {
	c.execute()
	app.history.push(&c)
}

// App can undo commands with command history stack
func (app *Applictation) undo() {
	command, ok := app.history.pop()
	if ok {
		command.undo()
	}
}

// Simplified text editor example struct
type Editor struct {
	selection string
}

func newEditor() *Editor {
	return &Editor{
		selection: "",
	}
}

func (e *Editor) getSelection() string {
	fmt.Println("Selected text is returned")
	return e.selection
}

func (e *Editor) deleteSelection() {
	fmt.Println("Selected text is deleted")
	e.selection = ""
}

func (e *Editor) replaceSelection(s string) {
	fmt.Println("Insert text from clipboard at selected position")
	e.selection = s
}

type CommandsList struct {
	Applictation
	Editor
	CopyCommand
	CutCommand
	PasteCommand
}

// Command stack for undo operations
type CommandHistory struct {
	commands []Command
}

func (ch *CommandHistory) push(c *Command) {
	ch.commands = append(ch.commands, *c)
}

func (ch *CommandHistory) pop() (com Command, isEmpty bool) {
	if len(ch.commands) > 0 {
		return nil, true
	}
	com = ch.commands[len(ch.commands)-1]
	ch.commands = ch.commands[:len(ch.commands)-1]
	return com, false
}

// interface has at least 1 func - execute
type Command interface {
	execute()
	undo()
	saveBackup()
}

type TemplateCommand struct { // this shall be embedded in concrete commands to reduce copies
	app    Applictation
	editor Editor
	backup string
}

func (tc *TemplateCommand) saveBackup() {
	tc.backup = tc.editor.selection
}

// concrete command 1
type CopyCommand struct {
	TemplateCommand
}

func NewCopyCommand(app *Applictation, e *Editor) *CopyCommand {
	return &CopyCommand{
		TemplateCommand: TemplateCommand{
			app:    *app,
			editor: *e,
			backup: "",
		},
	}
}

func (cp *CopyCommand) execute() {
	cp.app.clipboard = cp.editor.getSelection()
}

func (cp *CopyCommand) undo() {

}

type CutCommand struct {
	TemplateCommand
}

func NewCutCommand(app *Applictation, e *Editor) *CopyCommand {
	return &CopyCommand{
		TemplateCommand: TemplateCommand{
			app:    *app,
			editor: *e,
			backup: "",
		},
	}
}

func (ct *CutCommand) execute() {
	ct.saveBackup()
	ct.app.clipboard = ct.editor.getSelection()
	ct.editor.deleteSelection()
}

type PasteCommand struct {
	TemplateCommand
}

func NewPasteCommand(app *Applictation, e *Editor) *CopyCommand {
	return &CopyCommand{
		TemplateCommand: TemplateCommand{
			app:    *app,
			editor: *e,
			backup: "",
		},
	}
}

func (p *PasteCommand) execute() {
	p.saveBackup()
	p.editor.replaceSelection(p.app.clipboard)
}

type UndoCommand struct {
	TemplateCommand
}

func NewUndoCommand(app *Applictation, e *Editor) *CopyCommand {
	return &CopyCommand{
		TemplateCommand: TemplateCommand{
			app:    *app,
			editor: *e,
			backup: "",
		},
	}
}

func (uc *UndoCommand) undo() {
	uc.app.undo()
}

func testCommand() {
	// initialize editor & application
	editor := newEditor()
	app := NewApplication(editor)

	copy := NewCopyCommand(app, editor)
	paste := NewPasteCommand(app, editor)
	cut := NewCutCommand(app, editor)
	undo := NewUndoCommand(app, editor)
	// initialize buttons & link initial commands (they can be replaced if needed)
	copyButton := Button{command: copy}
	pasteButton := Button{command: paste}
	cutButton := Button{command: cut}
	undoButton := Button{command: undo}

	editor.selection = "initial"

	copyButton.command.execute()  // adding "initial" to app cliboard
	pasteButton.command.execute() // replacing initial with... initial
	cutButton.command.execute()   // cutting initial
	undoButton.command.execute()  // getting initial back (i hope so)

}
