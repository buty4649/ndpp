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
	"ndpp/ndp"
	"os/exec"
	"strings"
)

func commandExec(r []*ndp.Result, c string) error {
	var sb strings.Builder
	for _, v := range r {
		sb.WriteString("router_addr=")
		sb.WriteString(removeZoneIndex(v.Router.Addr.String()))
		sb.WriteString(" ")

		sb.WriteString("lladdr=")
		sb.WriteString(removeZoneIndex(v.Router.LLAddr))
		sb.WriteString(" ")

		sb.WriteString("local_addr=")
		sb.WriteString(removeZoneIndex(v.Local.Addr.String()))
		sb.WriteString(" ")

		sb.WriteString("interface=")
		sb.WriteString(v.IfName)
		sb.WriteString("\n")
	}

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
		i.Write([]byte(sb.String()))
	}()

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
