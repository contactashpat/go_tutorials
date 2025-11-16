package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"go_tutorials/internal/visualiser"
)

// SeeCommand encapsulates CLI parsing for the `see` sub-command.
type SeeCommand struct{}

// NewSeeCommand returns a ready-to-run SeeCommand.
func NewSeeCommand() *SeeCommand {
	return &SeeCommand{}
}

// Run executes the see command using provided CLI args.
func (c *SeeCommand) Run(args []string) error {
	fs := flag.NewFlagSet("see", flag.ContinueOnError)
	nameFlag := fs.String("name", "", "Name or tokens to visualize")
	reverseFlag := fs.String("reverse", "", "Reverse input: 'codepoints' or 'bytes'")
	if err := fs.Parse(args); err != nil {
		return err
	}

	input := strings.TrimSpace(*nameFlag)
	if input == "" {
		input = strings.Join(fs.Args(), " ")
	}
	if input == "" {
		return errors.New("no name provided; use --name or add it after the command")
	}

	resolved, note, err := c.resolveInput(*reverseFlag, input)
	if err != nil {
		return err
	}

	results, err := visualiser.AnalyseString(resolved)
	if err != nil {
		return err
	}

	renderTable(resolved, note, results)
	return nil
}

func (c *SeeCommand) resolveInput(reverseMode, input string) (string, string, error) {
	mode := strings.ToLower(strings.TrimSpace(reverseMode))
	if mode == "" {
		return input, "", nil
	}

	tokens := tokenizeReverseInput(input)
	if len(tokens) == 0 {
		return "", "", fmt.Errorf("reverse=%s requires values", reverseMode)
	}

	switch mode {
	case "codepoint", "codepoints", "cp":
		built, err := buildStringFromCodePoints(tokens)
		if err != nil {
			return "", "", err
		}
		return built, fmt.Sprintf("built from %s", reverseMode), nil
	case "byte", "bytes":
		built, err := buildStringFromBytes(tokens)
		if err != nil {
			return "", "", err
		}
		return built, fmt.Sprintf("built from %s", reverseMode), nil
	default:
		return "", "", fmt.Errorf("unknown reverse mode %q (use 'codepoints' or 'bytes')", reverseMode)
	}
}
