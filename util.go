package main

import (
	// "fmt"
	// "io/ioutil"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	// "time"
)

func run_command(name string, args []string, input string) (string, error) {

	var output bytes.Buffer
	c := exec.Command(name, args...)
	c.Stdout = &output
	stdin, err := c.StdinPipe()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", err
	}
	i, err := c.StderrPipe()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", err
	}

	if err = c.Start(); err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", err
	}
	if input != "" {
		stdin.Write([]byte(input))
		stdin.Close()
	}
	b, _ := ioutil.ReadAll(i)
	if err := c.Wait(); err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", err
	}
	fmt.Println(string(b))
	return string(b), nil

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
func run(id *Run) {
	log.Println("running" + id.Id)
	//mkdir workspace
	id.WorkDir = "/home/zpl/docker-run-test/" + id.Id + "/"
	err := add_work_path(id.WorkDir)
	if err != nil {
		// return http.StatusServiceUnavailable, Must(enc.Encode(
		// 	NewError(ErrRunCmd, fmt.Sprintf("server error in create workspace"))))
		log.Println("server error in create workspace")
		log.Println(err)
		return
	}
	log.Println("create dir " + id.WorkDir + " successful")
	err = os.Chdir(id.WorkDir)
	if err != nil {
		// return http.StatusServiceUnavailable, Must(enc.Encode(
		// 	NewError(ErrRunCmd, fmt.Sprintf("server error in create workspace"))))
		log.Println("server error in change  workspace")
		log.Println(err)
		return
	}
	log.Println("change dir " + id.WorkDir + " successful")
	err = put_code_into_file(id.WorkDir+id.Code.Filename, id.Code.Content)
	if err != nil {
		// return http.StatusServiceUnavailable, Must(enc.Encode(
		// 	NewError(ErrRunCmd, fmt.Sprintf("server error in create code file"))))
		log.Println("server error in create code file")
		log.Println(err)
		return
	}
	log.Println("render code ok")
	//handler with command
	for _, v := range id.Cmds {
		f, err := find_cmd_root(v.Cmd, id.WorkDir)
		if err != nil {
			// return http.StatusNotAcceptable, Must(enc.Encode(
			// 	NewError(ErrRunCmd, fmt.Sprintf("not found command %s", v.Cmd))))
			log.Println("not found command " + v.Cmd)
			return
		}
		res, err := run_command(f, v.Args, "")
		if err != nil {
			log.Println("error in run cmd" + v.Cmd)
			log.Println(err)
			return
		}

		log.Println("cmd " + v.Cmd + " run result :" + res)
	}
	// cmd := exec.Command("/bin/sh", "-c", "ping 127.0.0.1")
	// _, err := cmd.Output()
	// if err != nil {
	// 	panic(err.Error())
	// }

	// if err := cmd.Start(); err != nil {
	// 	panic(err.Error())
	// }

	// if err := cmd.Wait(); err != nil {
	// 	panic(err.Error())
	// }
}

// func main() {
// 	go run()
// 	time.Sleep(1e9)

// 	cmd := exec.Command("/bin/sh", "-c", `ps -ef | grep -v "grep" | grep "ping"`)
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		fmt.Println("StdoutPipe: " + err.Error())
// 		return
// 	}

// 	stderr, err := cmd.StderrPipe()
// 	if err != nil {
// 		fmt.Println("StderrPipe: ", err.Error())
// 		return
// 	}

// 	if err := cmd.Start(); err != nil {
// 		fmt.Println("Start: ", err.Error())
// 		return
// 	}

// 	bytesErr, err := ioutil.ReadAll(stderr)
// 	if err != nil {
// 		fmt.Println("ReadAll stderr: ", err.Error())
// 		return
// 	}

// 	if len(bytesErr) != 0 {
// 		fmt.Printf("stderr is not nil: %s", bytesErr)
// 		return
// 	}

// 	bytes, err := ioutil.ReadAll(stdout)
// 	if err != nil {
// 		fmt.Println("ReadAll stdout: ", err.Error())
// 		return
// 	}

// 	if err := cmd.Wait(); err != nil {
// 		fmt.Println("Wait: ", err.Error())
// 		return
// 	}

// 	fmt.Printf("stdout: %s", bytes)
// }
