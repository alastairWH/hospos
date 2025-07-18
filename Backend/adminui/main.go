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
	// Set window icon
	iconRes, err := fyne.LoadResourceFromPath("assets/icon-hospos.png")
	if err == nil {
		w.SetIcon(iconRes)
	}

	// Now define widgets after app and window are created
	refreshRolesBtn := widget.NewButton("Refresh Roles", func() {
		// Implement role refresh logic if needed
	})

	addRoleBtn := widget.NewButton("Add Role", nil) // OnTapped set below

	roleList := widget.NewMultiLineEntry()
	roleList.SetPlaceHolder("Roles will be listed here...")

	statusLabel := widget.NewLabel("Checking server...")
	go func() {
		for {
			online := adminapi.CheckAPIStatus()
			status := "Offline"
			if online {
				status = "Online"
			}
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

	// User management UI
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
		pinEntry := widget.NewEntry()
		pinEntry.SetPlaceHolder("3-6 digit PIN")
		roleSelect := widget.NewSelect([]string{}, nil)
		go func() {
			roles, err := adminapi.GetRoles()
			if err == nil {
				var roleNames []string
				for _, r := range roles {
					roleNames = append(roleNames, r.Role)
				}
				roleSelect.Options = roleNames
				roleSelect.Refresh()
			}
		}()
		dialog.ShowForm("Add User", "Add", "Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("Name", nameEntry),
				widget.NewFormItem("PIN", pinEntry),
				widget.NewFormItem("Role", roleSelect),
			},
			func(ok bool) {
				if ok {
					go func() {
						err := adminapi.AddUser(adminapi.User{Name: nameEntry.Text, Pin: pinEntry.Text, Role: roleSelect.Selected})
						if err != nil {
							dialog.ShowError(err, w)
						} else {
							refreshBtn.OnTapped()
						}
					}()
				}
			}, w)
	})

	addRoleBtn.OnTapped = func() {
		roleEntry := widget.NewEntry()
		dialog.ShowForm("Add Role", "Add", "Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("Role Name", roleEntry),
			},
			func(ok bool) {
				if ok {
					go func() {
						err := adminapi.AddRole(roleEntry.Text)
						if err != nil {
							dialog.ShowError(err, w)
						} else {
							refreshRolesBtn.OnTapped()
						}
					}()
				}
			}, w,
		)
	}

	// Devtools buttons (now in scope)
	dbInitBtn := widget.NewButton("DB Init", func() {
		go func() {
			err := adminapi.InitDB()
			if err != nil {
				fyne.CurrentApp().SendNotification(&fyne.Notification{Title: "DB Init", Content: "Error: " + err.Error()})
				dialog.ShowError(err, w)
			} else {
				fyne.CurrentApp().SendNotification(&fyne.Notification{Title: "DB Init", Content: "Database initialized!"})
				dialog.ShowInformation("DB Init", "Database initialized!", w)
			}
		}()
	})

	seedBtn := widget.NewButton("Seed Test Data", func() {
		go func() {
			err := adminapi.SeedTestData()
			if err != nil {
				fyne.CurrentApp().SendNotification(&fyne.Notification{Title: "Seed Data", Content: "Error: " + err.Error()})
				dialog.ShowError(err, w)
			} else {
				fyne.CurrentApp().SendNotification(&fyne.Notification{Title: "Seed Data", Content: "Test data seeded!"})
				dialog.ShowInformation("Seed Data", "Test data seeded!", w)
			}
		}()
	})

	clearBtn := widget.NewButton("Clear Test Data", func() {
		go func() {
			err := adminapi.ClearTestData()
			if err != nil {
				fyne.CurrentApp().SendNotification(&fyne.Notification{Title: "Clear Data", Content: "Error: " + err.Error()})
				dialog.ShowError(err, w)
			} else {
				fyne.CurrentApp().SendNotification(&fyne.Notification{Title: "Clear Data", Content: "Test data cleared!"})
				dialog.ShowInformation("Clear Data", "Test data cleared!", w)
			}
		}()
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("HOSPOS Backend Admin UI (Fyne)"),
		container.NewHBox(startBtn, stopBtn, statusLabel, dbInitBtn, seedBtn, clearBtn),
		widget.NewSeparator(),
		widget.NewLabel("User Management"),
		container.NewHBox(refreshBtn, addUserBtn),
		userList,
		widget.NewSeparator(),
		widget.NewLabel("Role Management"),
		container.NewHBox(refreshRolesBtn, addRoleBtn),
		roleList,
	))
	w.Resize(fyne.NewSize(500, 400))
	w.ShowAndRun()
}
