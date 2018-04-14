package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jinzhu/configor"
)

func main() {
	var Config = struct {
		Entries map[string]string
		Options map[string]string
	}{}

	configor.Load(&Config, "/home/jakob/.config/editconf.yaml")
	keys := make([]string, 0, len(Config.Entries))
	for key := range Config.Entries {
		keys = append(keys, key)
	}

	echo := exec.Command("echo", strings.Join(keys, "\n"))
	dmenu := exec.Command("dmenu", "-b")
	pipe, err := echo.StdoutPipe()
	if err != nil {
		panic(err)
	}
	dmenu.Stdin = pipe

	dmenu.Start()
	echo.Run()
	choice, er := dmenu.Output()
	dmenu.Wait()
	// dmenu.Run()

	fmt.Println(choice, er)
	if string(choice) != "" {
		path := Config.Entries[string(choice)]

		editor, ok := os.LookupEnv("EDITOR")
		fmt.Println(ok)
		edit := exec.Command(Config.Options["terminal"], "-e", editor, path)
		/*edit.Stdin = os.Stdin
		edit.Stdout = os.Stdout
		edit.Stderr = os.Stderr
		*/
		err := edit.Run()
		if err != nil {
			fmt.Println(err)
		}

		/*
			compiled, _ := regexp.MatchString("\\.h$", path)
			if compiled {
				cmd := exec.Command("make")
				cmd.Dir = path
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			}
		*/
	}

	/*
		for disp, path := range Config.Entries {
			fmt.Println(disp, path)
		}
	*/
	// fmt.Println("%#v", Config)
}
