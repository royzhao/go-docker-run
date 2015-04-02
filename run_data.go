package main

import (
	"encoding/xml"
)

type Code struct {
	filename string
	content  string
}
type Command struct {
	cmd  string
	args []string
}
type Run struct {
	XMLName xml.Name  `json:"-" xml:"run"`
	Id      string    `json:"id" xml:"id,attr"`
	WorkDir string    `json:"workdir" xml:"workdir"`
	Code    Code      `json:"code" xml:"code"`
	Cmds    []Command `json:"cmds" xml:"cmds"`
}
