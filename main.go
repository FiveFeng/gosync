package main

import (
	"os"
	"os/exec"
	"os/signal"

	"github.com/fivefeng/gosync/server"
)

func main() {
	// Gin模块
	chChromDie := make(chan struct{})
	chBackendDie := make(chan struct{})
	chSignal := listenToInterrupt()
	go server.Run()
	go startBrowser(chChromDie, chBackendDie)
	// lorca方式
	/* var ui lorca.UI
	ui, _ = lorca.New("http://localhost:8080", "", 800, 600, "--disable-sync", "--disable-translate")
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ui.Done():
	case <-chSignal:
	}
	ui.Close() */

	//cmd := startBrowser()
	// 监听ctrl+c
	for {
		select {
		case <-chSignal:
			chBackendDie <- struct{}{}
		case <-chChromDie:
			os.Exit(0)
		}
	}
}

func startBrowser(chChromDie chan struct{}, chBackendDie chan struct{}) {
	// exec方式
	//chromePath := "D:\\Google\\Chrome\\Application\\chrome.exe"
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://localhost:27149/static/index.html")
	cmd.Start()
	go func() {
		<-chBackendDie
		cmd.Process.Kill()
	}()
	go func() {
		cmd.Wait()
		//return cmd
		chChromDie <- struct{}{}
	}()
}

func listenToInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	return chSignal
}
