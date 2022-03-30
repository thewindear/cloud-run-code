package service

import (
	"bytes"
	context2 "context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var TimeoutError = errors.New("execute timeout")
var Success = errors.New("success")
var Exited = errors.New("exited")

var DockerRunner dockerRunner

func InitDockerRunner(config string) error {
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(config)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &DockerRunner)
	if err != nil {
		return err
	}
	DockerRunner.DockerPath = dockerPath
	return nil
}

type DRunner struct {
	Ext      string `json:"ext"`
	Filename string `json:"filename"`
	Image    string `json:"image"`
	Cmd      string `json:"cmd"`
}

type dockerRunner struct {
	Timeout    int                 `json:"timeout"`
	TmpPath    string              `json:"tmp_path"`
	DockerBase string              `json:"docker_base"`
	Runners    map[string]*DRunner `json:"docker_runner"`
	DockerPath string
}

func (dr dockerRunner) RunnerExists(runner string) bool {
	if _, ok := dr.Runners[runner]; ok {
		return true
	}
	return false
}

func (dr dockerRunner) Exec(ctx context2.Context, runnerName string, code string) (result []byte, err error) {
	runner, _ := dr.Runners[runnerName]
	tmpFileName := fmt.Sprintf("%d.%s_", time.Now().UnixMilli(), runner.Ext)
	var tmpFile *os.File
	if tmpFile, err = os.CreateTemp("", tmpFileName); err != nil {
		return nil, err
	}
	_, _ = tmpFile.WriteString(code)
	log.Println("tmp file path: " + tmpFile.Name())
	defer func() {
		_ = syscall.Unlink(tmpFile.Name())
	}()
	cmdStr := strings.Clone(dr.DockerBase)
	cmdStr = strings.Replace(cmdStr, "{tmp_file}", tmpFile.Name(), 1)
	cmdStr = strings.Replace(cmdStr, "{runner_filename}", runner.Filename, 1)
	cmdStr = strings.Replace(cmdStr, "{image}", runner.Image, 1)
	args := strings.Split(cmdStr, " ")
	args = append(args, runner.Cmd)
	result, err = dr.ExecDocker(ctx, args)
	return
}

func (dr dockerRunner) ExecDocker(ctx context2.Context, args []string) (result []byte, err error) {
	cmd := exec.CommandContext(ctx, dr.DockerPath, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err = cmd.Start(); err != nil {
		return
	}
	chErr := make(chan error)
	go func() {
		for {
			select {
			case <-ctx.Done():
				chErr <- TimeoutError
				_ = cmd.Process.Kill()
				return
			default:
				if state, err := cmd.Process.Wait(); err != nil {
					chErr <- err
					return
				} else {
					if state.Success() {
						chErr <- Success
						return
					}
					if state.Exited() {
						chErr <- Exited
						return
					}
				}
			}
		}
	}()

	err = <-chErr
	if err == Success || err == Exited {
		err = nil
	}
	if stderr.Len() > 0 {
		return stderr.Bytes(), err
	}
	return out.Bytes(), err
}
