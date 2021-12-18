package hackflow

import (
	"fmt"
	"os/exec"
)

type docker struct{}

func NewDocker() *docker {
	return &docker{}
}

func (d *docker) Install(link, dst string) (execPath string, err error) {
	err = exec.Command("docker", "pull", link).Run()
	return fmt.Sprintf("docker run -it  %s", link), err
}
