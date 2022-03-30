package test

import (
	"bytes"
	"cloud-run-code/service"
	"context"
	"errors"
	"log"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestLookUp(t *testing.T) {
	t.Log(exec.LookPath("docker"))
}

func TestRunner(t *testing.T) {
	_ = service.InitDockerRunner("/Users/edz/codeLab/cloud-run-code/config/docker_config.json")
	//content, err := dr.Exec("node", "const node1 = 'ABC'; console.log(`hello ${node1} node`);")
	javacode := `
public class App {
		public static void main(String[] args) {
			System.out.print("hello world 123");
		}
	}
`
	var start = time.Now().Unix()

	cancelCtx, cancelFn := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFn()

	content, err := service.DockerRunner.Exec(cancelCtx, "java", javacode)

	var end = time.Now().Unix()

	if err != nil {
		log.Println(err)
	} else {
		log.Println("output:" + string(content))
	}

	log.Println(end - start)

	//if err != nil {
	//	t.Fatalf(string(content))
	//} else {
	//	t.Log(string(content))
	//}
}

func TestPHPRunner(t *testing.T) {
	_ = service.InitDockerRunner("/Users/edz/codeLab/cloud-run-code/config/docker_config.json")
	//content, err := dr.Exec("node", "const node1 = 'ABC'; console.log(`hello ${node1} node`);")
	code := `
<?php
	sleep(3)
	echo 'aa';
`
	var start = time.Now().Unix()
	cancelCtx, _ := context.WithTimeout(context.Background(), time.Second*10)

	content, err := service.DockerRunner.Exec(cancelCtx, "php", code)

	var end = time.Now().Unix()

	if err != nil {
		log.Println(err)
	} else {
		log.Println("output:" + string(content))
	}

	log.Println(end - start)
}

func TestWithTimeout(t *testing.T) {
	_ = service.InitDockerRunner("/Users/edz/codeLab/cloud-run-code/config/docker_config.json")
	str := "run -i --cpus=1 -m=512M --rm --stop-timeout=1 --network none -v /var/folders/g0/mt1ssr8s2cj8k7z2hs1d8wmh0000gn/T/1648645701373.php_2810193720:/usr/src/myapp/app.php -w /usr/src/myapp php:7.4.28-zts-alpine sh -c"
	args := strings.Split(str, " ")
	args = append(args, "php app.php")

	cancelCtx, cancelFn := context.WithTimeout(context.Background(), time.Second*1)
	defer cancelFn()

	cmd := exec.CommandContext(cancelCtx, "docker", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	var err error
	if err = cmd.Start(); err != nil {
		log.Fatalf(err.Error())
	}
	ch := make(chan error)
	go func() {
		for {
			select {
			case <-cancelCtx.Done():
				log.Println("111")
				err = cmd.Process.Kill()
				ch <- errors.New("timeout:" + err.Error())
				return
			default:
				log.Println("def")
				state, err := cmd.Process.Wait()
				if err != nil {
					err = errors.Unwrap(err)
					ch <- err
					return
				}
				if state.Success() {
					ch <- errors.New("success")
					return
				}
				if state.Exited() {
					ch <- errors.New("exited")
					return
				}
			}
		}
	}()

	err = <-ch
	if err.Error() == "success" {
		log.Println(out.String())
	} else {
		log.Println(err.Error())
	}
}
