package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"net/http"

	"encoding/json"

	"bytes"
    "os"
	"github.com/alexellis/faas/gateway/requests"
)

func main() {
	// var handler string
	var source string
	var image string

	var action string
	var functionName string
	var gateway string

	flag.StringVar(&source, "source", "", "source file for function, i.e. handler.js")
	flag.StringVar(&image, "image", "", "Docker image name to build")
	flag.StringVar(&action, "action", "", "either build or deploy")
	flag.StringVar(&functionName, "name", "", "give the name of your deployed function")
	flag.StringVar(&gateway, "gateway", "http://localhost:8080", "gateway URI - i.e. http://localhost:8080")

	flag.Parse()

	if len(action) == 0 {
		fmt.Println("give either -action= build or deploy")
		return
	}

	if action == "build" {
		if len(image) == 0 {
			fmt.Println("Give a valid -image name for your Docker image.")
			return
		}

		fmt.Printf("Building: %s with Docker. Please wait..\n", image)

		builder := strings.Split(fmt.Sprintf("docker build --build-arg http_proxy=%s --build-arg https_proxy=%s -t %s .", os.Getenv("http_proxy"), os.Getenv("https_proxy"), image), " ")
        fmt.Println(strings.Join(builder, " "))
		targetCmd := exec.Command(builder[0], builder[1:]...)
		targetCmd.Dir = "./template/node/"
		cmdOutput, cmdErr := targetCmd.CombinedOutput()
		if cmdErr != nil {
			fmt.Println(cmdErr)
		}
		fmt.Println(string(cmdOutput))

		fmt.Printf("Image: %s built.\n", image)
	} else if action == "deploy" {
		if len(image) == 0 {
			fmt.Println("Give an image name to be deployed.")
			return
		}
		if len(functionName) == 0 {
			fmt.Println("Give a -name for your function as it will be deployed on FaaS")
			return
		}

		req := requests.CreateFunctionRequest{
			EnvProcess: "node index.js",
			Image:      image,
			Network:    "func_functions",
			Service:    functionName,
		}

		reqBytes, _ := json.Marshal(&req)
		reader := bytes.NewReader(reqBytes)
		res, err := http.Post(gateway+"/system/functions", "application/json", reader)
		if err != nil {
            fmt.Println("Is FaaS deployed? Do you need to specify the -gateway flag?")
			fmt.Println(err)
            return
		}
		fmt.Println(res.Status)
        deployedUrl := fmt.Sprintf("URL: %s/function/%s\n", gateway,functionName)
        fmt.Println(deployedUrl)

	}

}