package main

import (
	"fmt"
	"lillith/autorender/worker/config"
	"lillith/autorender/worker/portal2"
	"log"
	"net/http"
	"os"
	"time"
)

const FILE_PERM = 0755

func main() {
	err := config.ReadFile("worker_config.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Worker %s active\n", config.Cfg.Id)

	demoPath := fmt.Sprintf("%s/portal2/demos/autorender", config.Cfg.GameDir)

	_, err = os.Stat(demoPath)
	if err != nil {
		log.Printf("demos/autorender does not exist! Creating...\n")
		err = os.Mkdir(demoPath, FILE_PERM)
		if err != nil {
			log.Fatalln("Failed to create demos/autorender!")
		}
	}

	for {
		time.Sleep(time.Duration(config.Cfg.RequestInterval) * time.Second)

		request, err := http.NewRequest("GET", fmt.Sprintf("%s/api/fetch/queue", config.Cfg.Server), nil)
		if err != nil {
			log.Printf("Error creating request: %v\n", err)
			continue
		}

		request.Header.Add("Worker-ID", "sample")

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Printf("Error requesting demo: %v\n", err)
			continue
		}

		if response.Status == "200 OK" {
			log.Print("Claimed a demo\n\n")

			portal2.DeleteFiles()

			idbuf := make([]byte, response.ContentLength)
			response.Body.Read(idbuf)
			id := string(idbuf)

			err = portal2.DownloadDemo(id)
			if err != nil {
				log.Printf("error downloading demo: %v\n", err)
				continue
			}

			err = portal2.DownloadConfigs()
			if err != nil {
				log.Printf("Failed to download configs")
				continue
			}

			log.Println("Downloaded configs")

			portal2.PrepareAutoexec(id)

			game := portal2.Launch()
			log.Printf("Rendering demo %s\n", id)
			game.Wait()
			log.Printf("Finished rendering demo %s\n", id)

			err = portal2.UploadVideo(id)
			if err != nil {
				log.Printf("Failed to upload video %s\n", id)
			}

			log.Printf("Uploaded video %s\n", id)
		}
	}
}
