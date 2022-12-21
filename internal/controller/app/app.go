package app

import (
	"Statistic/internal/config"
	"Statistic/internal/controller"
	"Statistic/internal/storage"
	"context"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func GUIapp(s *controller.Server, wg *sync.WaitGroup, sig chan os.Signal) { //(*controller.Server, *sync.WaitGroup, *chan os.Signal, string, string) {
	wg.Add(1)
	var (
		listenPort, selfPort string
	)

	myApp := app.New()
	myWindowApp := myApp.NewWindow("config window")
	contentFirst := widget.NewLabel("specify simulator port\nDefault: port=\"8383\"")
	contentForEntry := widget.NewLabel("поле ввода")

	entryFirst := widget.NewEntry()
	buttonFirst := widget.NewButton("simulator\nport", func() {
		data := entryFirst.Text

		entryFirst.Hide()
		entryFirst.Disable()
		contentForEntry.Hide()
		contentForEntry.Hide()

		if len(data) == 4 {
			tmp, err := strconv.Atoi(data)
			if err != nil {
				log.Println("use default port for listen simulator")
				listenPort = config.ListenPort
			}
			log.Printf("use %d as listen port simulator", tmp)
			listenPort = data
		} else {
			log.Println("use default port for listen simulator")
			listenPort = config.ListenPort
		}
	})

	contentSecond := widget.NewLabel("specify show port\nDefault: port=\"8282\"")
	entrySecond := widget.NewEntry()
	buttonSecond := widget.NewButton("statistic\nport", func() {
		data := entrySecond.Text
		entrySecond.Hide()

		if len(data) == 4 {
			tmp, err := strconv.Atoi(data)
			if err != nil {
				log.Println("use default port for server statistics")
				selfPort = config.SelfPort
			}
			log.Printf("use %d as listen port simulator", tmp)
			selfPort = data
		} else {
			log.Println("use default port for server statistics")
			selfPort = config.SelfPort
		}

	})

	buttonStart := widget.NewButton("start", func() {
		buttonFirst.Hide()
		buttonFirst.Disable()
		buttonSecond.Hide()
		contentSecond.Hide()
		contentFirst.Hide()
		//myWindowApp.Hide()
		buttonSecond.Disable()
		s.MountHandlers()

		storage.ListenPort = listenPort
		server := &http.Server{Addr: ":" + selfPort, Handler: s.Router}
		open, err := url.Parse("http://localhost:" + selfPort)
		if err != nil {
			log.Println("err in generate url")
		}
		myApp.OpenURL(open)
		serverCtx, serverStopCtx := context.WithCancel(context.Background())

		go func() {
			defer wg.Done()
			<-sig
			// Shutdown signal with grace period of 10 seconds
			shutdownCtx, _ := context.WithTimeout(serverCtx, 10*time.Second)

			go func() {
				<-shutdownCtx.Done()
				if shutdownCtx.Err() == context.DeadlineExceeded {
					log.Fatal("graceful shutdown timed out.. forcing exit.")
				}
			}()

			// Trigger graceful shutdown
			err := server.Shutdown(shutdownCtx)
			if err != nil {
				log.Fatal(err)
			}
			serverStopCtx()
			myWindowApp.Close()

			myApp.Quit()

		}()

		// Run the server
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}

		// Wait for server context to be stopped
		<-serverCtx.Done()

	})

	buttonStop := widget.NewButton("stop", func() {
		sig <- syscall.SIGQUIT
		log.Print("button stop")
		myApp.Quit()
	})

	myWindowApp.SetContent(container.NewVBox(
		contentFirst,
		contentForEntry,
		entryFirst,
		buttonFirst,

		contentSecond,
		entrySecond,
		buttonSecond,
		buttonStart,
		buttonStop,
	))

	myWindowApp.ShowAndRun()

}
