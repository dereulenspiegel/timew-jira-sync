package timewarrior

import (
	"fmt"
	"os"
	"os/exec"
)

// Cli gives access to the timewarrior CLI to interact with timewarrior
var Cli = Command{}

type Command struct{}

// Tag tags a specific interval with one or multiple tags
func (c Command) Tag(id uint64, tags ...string) error {
	args := append([]string{"tag", fmt.Sprintf("@%d", id)}, tags...)
	cmd := exec.Command("timew", args...)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to tag interval: %w (%s)", err, string(out))
	}
	return nil
}
