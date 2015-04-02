package main

import (
	"encoding/xml"
)

type Code struct {
	Filename string
	Content  string
}
type Command struct {
	Cmd  string
	Args []string
}
type Run struct {
	XMLName xml.Name  `json:"-" xml:"run"`
	Id      string    `json:"id" xml:"id,attr"`
	WorkDir string    `json:"workdir" xml:"workdir"`
	Code    Code      `json:"code" xml:"code"`
	Cmds    []Command `json:"cmds" xml:"cmds"`
}
