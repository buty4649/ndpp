/*
Copyright © 2024 buty4649

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package cmd

import (
	"bytes"
	"fmt"
	"ndpp/ndp"
	"os/exec"
)

func commandExec(r []*ndp.Result, c string) error {
	cmd := exec.Command(c)

	var o, e bytes.Buffer
	cmd.Stdout = &o
	cmd.Stderr = &e

	i, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer i.Close()

		vars, _ := formatResult(r, "shellvar")
		i.Write([]byte(vars))
	}()

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute command: %s, stdout: %s, stderr: %s", err, o.String(), e.String())
	}

	return nil
}
