package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func run_command(name string, args []string, input string) (string, error) {

	log.Println("run command" + name)
	c := exec.Command(name, args...)
	e, err := c.StderrPipe()
	if err != nil {
		return "error start cmd: " + name, err
	}
	o, err := c.StdoutPipe()
	if err != nil {
		return "error start cmd: " + name, err
	}
	err = c.Start()
	if err != nil {
		return "error start cmd: " + name, err
	}
	//start a timer,limit 5s
	timer := time.NewTicker(time.Second * 5)
	go func() {
		<-timer.C
		_, err := c.Stderr.Write([]byte("run too long"))
		if err != nil {
			log.Println(err)
		}
		c.Process.Signal(os.Kill)
		log.Println("cancel run command:" + name)
	}()
	ebyte, _ := ioutil.ReadAll(e)
	obyte, _ := ioutil.ReadAll(o)
	err = c.Wait()
	//cancel timer
	timer.Stop()
	log.Println("run command over")
	if err != nil {
		return string(ebyte), err
	}

	return string(obyte), nil

}
func put_code_into_file(filename string, code string) error {
	fout, err := os.Create(filename)
	defer fout.Close()
	if err != nil {
		fmt.Println(filename, err)
		return err
	}
	fout.WriteString(code)
	return nil
}
func add_work_path(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	return err
}
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
func find_cmd_root(cmd string, wd string) (string, error) {
	f, err := exec.LookPath(cmd)
	if err != nil {
		if Exist(wd+cmd) == true {
			return wd + cmd, nil
		}
		fmt.Println("not install " + cmd)
		return "", err
	}
	return f, nil
}
func CleanDir(path string) error {
	err := os.RemoveAll(path)
	return err
}
func run(id *Run) {
	var total_res = "no output\n"
	log.Println("running" + id.Id)
	//mkdir workspace
	id.WorkDir = "/home/zpl/docker-run-test/" + id.Id + "/"
	err := add_work_path(id.WorkDir)
	if err != nil {
		log.Println("server error in create workspace")
		log.Println(err)
		redis_client.Set(id.Id, []byte("server error in create workspace :"+err.Error()))
		return
	}
	log.Println("create dir " + id.WorkDir + " successful")
	err = os.Chdir(id.WorkDir)
	if err != nil {
		log.Println("server error in change  workspace")
		log.Println(err)
		redis_client.Set(id.Id, []byte("server error in change  workspace :"+err.Error()))
		return
	}
	log.Println("change dir " + id.WorkDir + " successful")
	err = put_code_into_file(id.WorkDir+id.Code.Filename, id.Code.Content)
	if err != nil {
		log.Println("server error in create code file")
		log.Println(err)
		redis_client.Set(id.Id, []byte("server error in create code file :"+err.Error()))
		return
	}
	log.Println("render code ok")
	//handler with command
	for _, v := range id.Cmds {
		f, err := find_cmd_root(v.Cmd, id.WorkDir)
		if err != nil {
			log.Println("not found command " + v.Cmd)
			redis_client.Set(id.Id, []byte("not found command "+v.Cmd))
			return
		}
		res, err := run_command(f, v.Args, "")
		if err != nil {
			log.Println("==================")
			log.Println("run cmd " + v.Cmd + " error")
			log.Println(res)
			log.Println(err)
			log.Println("******************")
			redis_client.Set(id.Id, []byte("run cmd "+v.Cmd+" error\n"+res+"\n"+err.Error()))
			return
		}
		log.Println("cmd " + v.Cmd + " run result :" + res)
		total_res += "cmd " + v.Cmd + " run result is :" + res + "\n"
	}
	log.Println(total_res)
	redis_client.Set(id.Id, []byte(total_res))
	//clean
	err = CleanDir(id.WorkDir)
	if err != nil {
		log.Println("Clean workspace is failed!")
	}
}
