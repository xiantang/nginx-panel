package nginx

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

var filePath = "/etc/nginx/tests/test.conf"

// Test do nginx test nginx -t
func Test(body string) error {
	if os.Getenv("SKIP_NGINX_TEST") == "true" {
		return nil
	}
	logrus.WithField("config", body).Info("nginx test")
	// write file to /etc/nginx/test/ than do nginx -t

	if err := ioutil.WriteFile(filePath, []byte(body), 0644); err != nil {
		logrus.WithField("config", body).WithError(err).Error("write file error")
		return err
	}
	defer os.Remove(filePath)
	cmd := exec.Command("nginx", "-c", "/etc/nginx/nginx_test.conf", "-t")
	buf := bytes.NewBufferString("")
	cmd.Stderr = buf
	err := cmd.Run()
	if err != nil {
		logrus.WithField("output", buf.String()).WithError(err).Error("nginx -t error")
		return err
	}

	return nil
}
