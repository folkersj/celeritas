package main

import (
    "fmt"
    "github.com/fatih/color"
    "github.com/go-git/go-git/v5"
    "io"
    "os"
    "os/exec"
    "runtime"
    "strings"
)

var appURL string

func doNew(appName string) {
    appName = strings.ToLower(appName)
    appURL = appName

    if strings.Contains(appName, "/") {
        exploded := strings.SplitAfter(appName, "/")
        appName = exploded[len(exploded)-1]
    }

    color.Green("\tCloning repository...")
    _, err := git.PlainClone("./"+appName, false, &git.CloneOptions{
        URL:      "https://github.com/folkersj/celeritas-app.git",
        Progress: os.Stdout,
        Depth:    1,
    })
    if err != nil {
        exitGracefully(err)
    }

    err = os.RemoveAll(fmt.Sprintf("./%s/.git", appName))
    if err != nil {
        exitGracefully(err)
    }

    color.Yellow("\tCreating .env file...")
    data, err := templateFS.ReadFile("templates/env.txt")
    if err != nil {
        exitGracefully(err)
    }
    env := string(data)
    env = strings.ReplaceAll(env, "${APP_NAME}", appName)
    env = strings.ReplaceAll(env, "${KEY}", cel.RandomString(32))
    err = copyDataToFile([]byte(env), fmt.Sprintf("./%s/.env", appName))
    if err != nil {
        exitGracefully(err)
    }

    if runtime.GOOS == "windows" {
        source, err := os.Open(fmt.Sprintf("./%s/Makefile.windows", appName))
        if err != nil {
            exitGracefully(err)
        }
        defer source.Close()

        destination, err := os.Create(fmt.Sprintf("./%s/Makefile", appName))
        if err != nil {
            exitGracefully(err)
        }
        defer destination.Close()

        _, err = io.Copy(destination, source)
        if err != nil {
            exitGracefully(err)
        }
    } else {
        source, err := os.Open(fmt.Sprintf("./%s/Makefile.mac", appName))
        if err != nil {
            exitGracefully(err)
        }
        defer source.Close()

        destination, err := os.Create(fmt.Sprintf("./%s/Makefile", appName))
        if err != nil {
            exitGracefully(err)
        }
        defer destination.Close()

        _, err = io.Copy(destination, source)
        if err != nil {
            exitGracefully(err)
        }
    }
    _ = os.Remove("./" + appName + "/Makefile.mac")
    _ = os.Remove("./" + appName + "/Makefile.windows")

    color.Yellow("\tCreating the go.mod file")
    _ = os.Remove("./" + appName + "/go.mod")
    data, err = templateFS.ReadFile("templates/go.mod.txt")
    if err != nil {
        exitGracefully(err)
    }
    mod := string(data)
    mod = strings.ReplaceAll(mod, "${APP_NAME}", appURL)
    err = copyDataToFile([]byte(mod), "./"+appName+"/go.mod")
    if err != nil {
        exitGracefully(err)
    }

    _ = os.Chdir("./" + appName)
    color.Yellow("\tUpdating source files...")
    updateSource()

    color.Yellow("\tRunning go mod tidy...")
    cmd := exec.Command("go", "mod", "tidy")
    err = cmd.Start()
    if err != nil {
        exitGracefully(err)
    }

    color.Green("Done building " + appURL)
    color.Green("Go build something awesome")
}
