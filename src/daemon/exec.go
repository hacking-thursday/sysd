// build linux

package daemon

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"daemon/execdriver"
	"daemon/execdriver/lxc"
	"github.com/docker/docker/engine"
	"github.com/docker/docker/pkg/broadcastwriter"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/log"
	"github.com/docker/docker/pkg/promise"
	"github.com/docker/docker/runconfig"
	"github.com/docker/docker/utils"
)

type execConfig struct {
	sync.Mutex
	ID            string
	Running       bool
	ProcessConfig execdriver.ProcessConfig
	StreamConfig
	OpenStdin  bool
	OpenStderr bool
	OpenStdout bool
	Container  *Container
}

type execStore struct {
	s map[string]*execConfig
	sync.Mutex
}

func newExecStore() *execStore {
	return &execStore{s: make(map[string]*execConfig, 0)}
}

func (e *execStore) Add(id string, execConfig *execConfig) {
	e.Lock()
	e.s[id] = execConfig
	e.Unlock()
}

func (e *execStore) Get(id string) *execConfig {
	e.Lock()
	res := e.s[id]
	e.Unlock()
	return res
}

func (e *execStore) Delete(id string) {
	e.Lock()
	delete(e.s, id)
	e.Unlock()
}

func (execConfig *execConfig) Resize(h, w int) error {
	return execConfig.ProcessConfig.Terminal.Resize(h, w)
}

func (d *Daemon) registerExecCommand(execConfig *execConfig) {
	// Storing execs in container inorder to kill them gracefully whenever the container is stopped or removed.
	execConfig.Container.execCommands.Add(execConfig.ID, execConfig)
	// Storing execs in daemon for easy access via remote API.
	d.execCommands.Add(execConfig.ID, execConfig)
}

func (d *Daemon) getExecConfig(name string) (*execConfig, error) {
	if execConfig := d.execCommands.Get(name); execConfig != nil {
		if !execConfig.Container.IsRunning() {
			return nil, fmt.Errorf("Container %s is not running", execConfig.Container.ID)
		}
		return execConfig, nil
	}

	return nil, fmt.Errorf("No such exec instance '%s' found in daemon", name)
}

func (d *Daemon) unregisterExecCommand(execConfig *execConfig) {
	execConfig.Container.execCommands.Delete(execConfig.ID)
	d.execCommands.Delete(execConfig.ID)
}

func (d *Daemon) getActiveContainer(name string) (*Container, error) {
	container := d.Get(name)

	if container == nil {
		return nil, fmt.Errorf("No such container: %s", name)
	}

	if !container.IsRunning() {
		return nil, fmt.Errorf("Container %s is not running", name)
	}

	return container, nil
}

func (d *Daemon) ContainerExecCreate(job *engine.Job) engine.Status {
	if len(job.Args) != 1 {
		return job.Errorf("Usage: %s [options] container command [args]", job.Name)
	}

	if strings.HasPrefix(d.execDriver.Name(), lxc.DriverName) {
		return job.Error(lxc.ErrExec)
	}

	var name = job.Args[0]

	container, err := d.getActiveContainer(name)
	if err != nil {
		return job.Error(err)
	}

	config := runconfig.ExecConfigFromJob(job)

	entrypoint, args := d.getEntrypointAndArgs(nil, config.Cmd)

	processConfig := execdriver.ProcessConfig{
		Privileged: config.Privileged,
		User:       config.User,
		Tty:        config.Tty,
		Entrypoint: entrypoint,
		Arguments:  args,
	}

	execConfig := &execConfig{
		ID:            utils.GenerateRandomID(),
		OpenStdin:     config.AttachStdin,
		OpenStdout:    config.AttachStdout,
		OpenStderr:    config.AttachStderr,
		StreamConfig:  StreamConfig{},
		ProcessConfig: processConfig,
		Container:     container,
		Running:       false,
	}

	d.registerExecCommand(execConfig)

	job.Printf("%s\n", execConfig.ID)

	return engine.StatusOK
}

func (d *Daemon) ContainerExecStart(job *engine.Job) engine.Status {
	if len(job.Args) != 1 {
		return job.Errorf("Usage: %s [options] exec", job.Name)
	}

	var (
		cStdin           io.ReadCloser
		cStdout, cStderr io.Writer
		cStdinCloser     io.Closer
		execName         = job.Args[0]
	)

	execConfig, err := d.getExecConfig(execName)
	if err != nil {
		return job.Error(err)
	}

	func() {
		execConfig.Lock()
		defer execConfig.Unlock()
		if execConfig.Running {
			err = fmt.Errorf("Error: Exec command %s is already running", execName)
		}
		execConfig.Running = true
	}()
	if err != nil {
		return job.Error(err)
	}

	log.Debugf("starting exec command %s in container %s", execConfig.ID, execConfig.Container.ID)
	container := execConfig.Container

	if execConfig.OpenStdin {
		r, w := io.Pipe()
		go func() {
			defer w.Close()
			io.Copy(w, job.Stdin)
		}()
		cStdin = r
		cStdinCloser = job.Stdin
	}
	if execConfig.OpenStdout {
		cStdout = job.Stdout
	}
	if execConfig.OpenStderr {
		cStderr = job.Stderr
	}

	execConfig.StreamConfig.stderr = broadcastwriter.New()
	execConfig.StreamConfig.stdout = broadcastwriter.New()
	// Attach to stdin
	if execConfig.OpenStdin {
		execConfig.StreamConfig.stdin, execConfig.StreamConfig.stdinPipe = io.Pipe()
	} else {
		execConfig.StreamConfig.stdinPipe = ioutils.NopWriteCloser(ioutil.Discard) // Silently drop stdin
	}

	attachErr := d.Attach(&execConfig.StreamConfig, execConfig.OpenStdin, false, execConfig.ProcessConfig.Tty, cStdin, cStdinCloser, cStdout, cStderr)

	execErr := make(chan error)

	// Remove exec from daemon and container.
	defer d.unregisterExecCommand(execConfig)

	go func() {
		err := container.Exec(execConfig)
		if err != nil {
			execErr <- fmt.Errorf("Cannot run exec command %s in container %s: %s", execName, container.ID, err)
		}
	}()

	select {
	case err := <-attachErr:
		if err != nil {
			return job.Errorf("attach failed with error: %s", err)
		}
		break
	case err := <-execErr:
		return job.Error(err)
	}

	return engine.StatusOK
}

func (d *Daemon) Exec(c *Container, execConfig *execConfig, pipes *execdriver.Pipes, startCallback execdriver.StartCallback) (int, error) {
	return d.execDriver.Exec(c.command, &execConfig.ProcessConfig, pipes, startCallback)
}

func (container *Container) Exec(execConfig *execConfig) error {
	container.Lock()
	defer container.Unlock()

	waitStart := make(chan struct{})

	callback := func(processConfig *execdriver.ProcessConfig, pid int) {
		if processConfig.Tty {
			// The callback is called after the process Start()
			// so we are in the parent process. In TTY mode, stdin/out/err is the PtySlave
			// which we close here.
			if c, ok := processConfig.Stdout.(io.Closer); ok {
				c.Close()
			}
		}
		close(waitStart)
	}

	// We use a callback here instead of a goroutine and an chan for
	// syncronization purposes
	cErr := promise.Go(func() error { return container.monitorExec(execConfig, callback) })

	// Exec should not return until the process is actually running
	select {
	case <-waitStart:
	case err := <-cErr:
		return err
	}

	return nil
}

func (container *Container) monitorExec(execConfig *execConfig, callback execdriver.StartCallback) error {
	var (
		err      error
		exitCode int
	)

	pipes := execdriver.NewPipes(execConfig.StreamConfig.stdin, execConfig.StreamConfig.stdout, execConfig.StreamConfig.stderr, execConfig.OpenStdin)
	exitCode, err = container.daemon.Exec(container, execConfig, pipes, callback)
	if err != nil {
		log.Errorf("Error running command in existing container %s: %s", container.ID, err)
	}

	log.Debugf("Exec task in container %s exited with code %d", container.ID, exitCode)
	if execConfig.OpenStdin {
		if err := execConfig.StreamConfig.stdin.Close(); err != nil {
			log.Errorf("Error closing stdin while running in %s: %s", container.ID, err)
		}
	}
	if err := execConfig.StreamConfig.stdout.Clean(); err != nil {
		log.Errorf("Error closing stdout while running in %s: %s", container.ID, err)
	}
	if err := execConfig.StreamConfig.stderr.Clean(); err != nil {
		log.Errorf("Error closing stderr while running in %s: %s", container.ID, err)
	}
	if execConfig.ProcessConfig.Terminal != nil {
		if err := execConfig.ProcessConfig.Terminal.Close(); err != nil {
			log.Errorf("Error closing terminal while running in container %s: %s", container.ID, err)
		}
	}

	return err
}
