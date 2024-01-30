package cmd_chain

import (
	"fmt"
	"strings"
)

// This is actually an example of command chain, not cmd-pipe. But i could't make up anything better yet
type CMD struct {
	inputData string
	cmdType   string
	args      []string
	result    string
	nextCMD   *CMD
}

type iCMD interface {
	Execute([]string) string
}

func (c *CMD) Execute(inputData string) {
	if inputData != "" {
		c.inputData = inputData
	}
	execCommand(c)
}

func parseArgs(input string) []string {
	cmds := strings.Split(input, "|")
	return cmds
}

func formCommandPipe(cmds []CMD, bodies []string) {
	for i := 0; i < len(cmds); i++ {
		parts := strings.Fields(bodies[i])
		cmds[i].cmdType = parts[0]
		cmds[i].args = parts[1:]
	}
	if len(cmds) > 1 {
		for i := 0; i < len(cmds)-1; i++ {
			cmds[i].nextCMD = &cmds[i+1]
		}
	}

}

func parseCMDArgs(body string) []string {
	return strings.Fields(body)
}

func handleInput() {
	var s string
	fmt.Scan(&s)
	cmdb := parseArgs(s)
	cmds := make([]CMD, len(cmdb))
	formCommandPipe(cmds, cmdb)
	// Start the chain
	cmds[0].Execute(cmds[0].inputData)
}

func parseArgsWithoutCMDName(s string) []string {
	return strings.Fields(s)
}

func execCommand(c *CMD) {

	result := ""
	if c.inputData != "" {
		c.args[len(c.args)-1] = c.inputData
	}

	switch c.cmdType {
	case "cd":
		c.result = CD(c.args)
		//... other commands
	default:
		fmt.Println("Usage:...")
		return
	}

	if result != "" {
		c.nextCMD.inputData = result
	}
	c.nextCMD.Execute(c.result)
}

func CD(args []string) string {
	// CD...
	return "result"
}
