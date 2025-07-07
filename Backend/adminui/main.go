package main

import (
	"fmt"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"hospos-backend/adminui/adminapi"
)

var serverCmd *exec.Cmd

func main() {
	a := app.New()
	w := a.NewWindow("HOSPOS Backend Admin")

	statusLabel := widget.NewLabel("Checking server...")
	go func() {
		for {
			online := adminapi.CheckAPIStatus()
			status := "Offline"
			if online {
				status = "Online"
			}
			fyne.CurrentApp().SendNotification(&fyne.Notification{Title: "Status", Content: "Server status: " + status})
			statusLabel.SetText("Server status: " + status)
			time.Sleep(3 * time.Second)
		}
	}()

	startBtn := widget.NewButton("Start Server", func() {
		if serverCmd == nil {
			serverCmd = exec.Command("go", "run", "main.go")
			go func() {
				err := serverCmd.Start()
				if err != nil {
					statusLabel.SetText("Failed to start: " + err.Error())
					serverCmd = nil
					return
				}
				statusLabel.SetText("Server running (PID: " + fmt.Sprint(serverCmd.Process.Pid) + ")")
				serverCmd.Wait()
				statusLabel.SetText("Server stopped")
				serverCmd = nil
			}()
		}
	})
	stopBtn := widget.NewButton("Stop Server", func() {
		if serverCmd != nil && serverCmd.Process != nil {
			_ = serverCmd.Process.Kill()
			statusLabel.SetText("Server stopped")
			serverCmd = nil
		}
	})

	userList := widget.NewMultiLineEntry()
	userList.SetPlaceHolder("Users will be listed here...")
	refreshBtn := widget.NewButton("Refresh Users", func() {
		go func() {
			users, err := adminapi.GetUsers()
			if err != nil {
				userList.SetText("Error: " + err.Error())
				return
			}
			var lines string
			for _, u := range users {
				lines += fmt.Sprintf("Name: %s, Role: %s\n", u.Name, u.Role)
			}
			userList.SetText(lines)
		}()
	})
	addUserBtn := widget.NewButton("Add User", func() {
		nameEntry := widget.NewEntry()
		roleEntry := widget.NewEntry()
		dialog.ShowForm("Add User", "Add", "Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("Name", nameEntry),
				widget.NewFormItem("Role", roleEntry),
			},
			func(ok bool) {
				if ok {
					go func() {
						err := adminapi.AddUser(adminapi.User{Name: nameEntry.Text, Role: roleEntry.Text})
						if err != nil {
							dialog.ShowError(err, w)
						} else {
							refreshBtn.OnTapped()
						}
					}()
				}
			}, w)
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("HOSPOS Backend Admin UI (Fyne)"),
		container.NewHBox(startBtn, stopBtn, statusLabel),
		widget.NewSeparator(),
		widget.NewLabel("User Management"),
		container.NewHBox(refreshBtn, addUserBtn),
		userList,
	))
	w.Resize(fyne.NewSize(500, 400))
	w.ShowAndRun()
}
