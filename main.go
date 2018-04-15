package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/jinzhu/configor"
)

func main() {
	var Config = struct {
		Entries map[string]string
		Options map[string]string
		Dmenu   map[string]string
	}{}

	confFile, _ := expand("~/.config/editconf.yaml")
	configor.Load(&Config, confFile)

	keys := make([]string, 0, len(Config.Entries))
	for key := range Config.Entries {
		keys = append(keys, key)
	}

	flags := ""
	for k, v := range Config.Dmenu {
		if string(v) == "true" {
			flags = fmt.Sprintf("%s -%s", flags, k)
		} else {
			flags = fmt.Sprintf("%s -%s '%s'", flags, k, v)
		}
	}
	// fmt.Println(flags)
	// strcmd := fmt.Sprintf("echo %s | tr ' ' '\\n' | dmenu %s -fn 'fira mono:12' -p 'config:'", strings.Join(keys, " "), bflag)
	strcmd := fmt.Sprintf("echo %s | tr ' ' '\\n' | dmenu%s", strings.Join(keys, " "), flags)
	// fmt.Println(strcmd)
	shell := os.Getenv("SHELL")
	choiceBytes, err := exec.Command(shell, "-c", strcmd).Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	choice := strings.TrimSpace(string(choiceBytes))
	unexpandedFile, ok := Config.Entries[choice]

	if choice != "" && ok {
		file, er := expand(unexpandedFile)
		if er != nil {
			fmt.Println("error expanding ~ to /home/$USER")
			return
		}
		dir := path.Dir(file)

		editor := Config.Options["editor"]
		if editor == "" {
			editor, _ = os.LookupEnv("EDITOR") //fallback
			if editor == "" {
				fmt.Println("err: $EDITOR is not defined and isn't set in config")
			}
		}
		// strcmd := fmt.Sprintf("%s -e %s %s", Config.Options["terminal"], editor, file)
		// edit := exec.Command(shell, "-c", strcmd) //using shell to expand ~/ to /home/user
		edit := exec.Command(Config.Options["terminal"], "-e", editor, file)
		if err := edit.Run(); err != nil {
			fmt.Println("error running ", Config.Options["terminal"], " with $EDITOR: ", err)
			return
		}

		_, err := os.Stat(dir + "/Makefile") //check for Makefile...
		m1 := os.IsExist(err)
		_, err = os.Stat(dir + "/makefile") //and makefile lower case...
		m2 := os.IsExist(err)

		//for the moment, this executes before vim is closed, which is essentially useless
		if m1 || m2 {
			if _, err := os.Stat(dir + "/makefile"); os.IsExist(err) {
				cmd := exec.Command("make", choice)
				cmd.Dir = dir
				s, e := cmd.CombinedOutput()
				if e != nil {
					fmt.Println("error running make:", err)
					return
				}
				fmt.Print(string(s))
			}
		}
	}
}

func expand(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}
