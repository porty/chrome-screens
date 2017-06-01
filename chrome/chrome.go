package chrome

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

type Chrome struct {
	executable string
	url        string
	process    *os.Process
}

func New() (*Chrome, error) {
	path, err := exec.LookPath("chromium-browser")
	if err != nil {
		path, err = exec.LookPath("google-chrome")
		if err != nil {
			return nil, err
		}
	}
	return &Chrome{
		executable: path,
	}, nil
}

func (c *Chrome) SetURL(url string) {
	c.url = url
	if err := c.stop(); err != nil {
		log.Print("Received error from stop: " + err.Error())
	}
	if err := c.start(); err != nil {
		log.Print("Received error from start: " + err.Error())
	}
}

func (c *Chrome) GetURL() string {
	return c.url
}

func (c *Chrome) stop() error {
	if c.process != nil {
		if err := c.process.Signal(syscall.SIGTERM); err != nil {
			return err
		}
		// if err := c.process.Kill(); err != nil {
		// 	return err
		// }
		// c.process.Kill()
		c.process = nil
	}
	return nil
}

func (c *Chrome) start() error {

	cmd := exec.Command(c.executable, "--disable-infobars", "--disable-translate", "--kiosk", "--temp-profile", "--password-store=basic", c.url)
	writer := PrefixWriter{
		prefix: "CHROME: ",
		writer: os.Stdout,
	}
	cmd.Stdout = &writer
	cmd.Stderr = &writer
	env := os.Environ()
	env = append(env, "DISPLAY=:0")
	cmd.Env = env
	if err := cmd.Start(); err != nil {
		return err
	}
	c.process = cmd.Process
	return nil
}
