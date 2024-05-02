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
	"encoding/json"
	"fmt"
	"io"
	"ndpp/ndp"
	"os"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v3"
)

type formattedRouterInfo struct {
	Addr          string `json:"addr" yaml:"addr"`
	AddrWithZone  string `json:"addr_with_zone" yaml:"addr_with_zone"`
	LinkLayerAddr string `json:"lladdr" yaml:"lladdr"`
}

type formattedLocalInfo struct {
	Addr          string `json:"addr" yaml:"addr"`
	Interface     string `json:"interface" yaml:"interface"`
	LinkLayerAddr string `json:"lladdr" yaml:"lladdr"`
}

type formattedResult struct {
	Local  formattedLocalInfo  `json:"local" yaml:"local"`
	Router formattedRouterInfo `json:"router" yaml:"router"`
}

type formattedResults []formattedResult

func outputResults(r []*ndp.Result, t, o string) error {
	f, err := formatResult(r, t)
	if err != nil {
		return err
	}

	var output io.Writer
	if o == "" {
		output = os.Stdout
	} else {
		f, err := os.Create(o)
		if err != nil {
			return err
		}
		defer f.Close()
		output = f
	}
	fmt.Fprint(output, f)
	return err
}

func formatResult(r []*ndp.Result, t string) (string, error) {
	var results formattedResults
	for _, v := range r {
		ra := v.Router.Addr.String()
		results = append(results, formattedResult{
			Local: formattedLocalInfo{
				Addr:          removeZoneIndex(v.Local.Addr.String()),
				Interface:     v.IfName,
				LinkLayerAddr: v.Local.LLAddr,
			},
			Router: formattedRouterInfo{
				Addr:          removeZoneIndex(ra),
				AddrWithZone:  ra,
				LinkLayerAddr: v.Router.LLAddr,
			},
		})
	}

	var output string
	var err error
	if t == "yaml" {
		output, err = results.yamlDecode()
	} else if t == "json" {
		output, err = results.jsonDecode()
	} else if t == "bird" {
		output, err = results.birdConfigDecode()
	}

	return output, err
}

func removeZoneIndex(a string) string {
	return strings.Split(a, "%")[0]
}

func (r formattedResults) yamlDecode() (string, error) {
	b, err := yaml.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r formattedResults) jsonDecode() (string, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r formattedResults) birdConfigDecode() (string, error) {
	const tmpl = `{{- range . -}}
define {{.Local.Interface | replaceDot}}_neighbor_address = {{.Router.Addr}};
define {{.Local.Interface | replaceDot}}_local_address = {{.Local.Addr}};
define {{.Local.Interface | replaceDot}}_lladdr = hex:{{.Local.LinkLayerAddr}};

#
# The following is a sample configuration. If you wish to use it, please uncomment the settings below and paste them into bird.conf.
#
#define remote_asn = <Please set remote ASN>;
#define my_asn = <Please set local ASN>;
#
#protocol radv {
#  interface "{{.Local.Interface}}" {
#    custom option type 1 value uplink_lladdr;
#  };
#}
#protocol bgp uplink {
#  neighbor {{.Local.Interface | replaceDot}}_neighbor_address as remote_asn;
#  interface "{{.Local.Interface}}";
#  local as my_asn;
#  ipv4 { extended next hop; import all; export all; };
#  ipv6 { import all; export all; };
#}

{{ end }}`

	t, err := template.New("bird").Funcs(template.FuncMap{
		"replaceDot": func(s string) string { return strings.ReplaceAll(s, ".", "_") },
	}).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, r)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
