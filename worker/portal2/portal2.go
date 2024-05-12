package portal2

import (
	"bytes"
	"fmt"
	"io"
	"lillith/autorender/worker/config"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func Launch() *exec.Cmd {
	game := exec.Command(fmt.Sprintf("%s/portal2.exe", config.Cfg.GameDir), "-windowed", "-novid")
	game.Start()
	return game
}

func DownloadDemo(id string) error {
	downloadResponse, err := http.Get(fmt.Sprintf("%s/storage/demos/%s.dem", config.Cfg.Server, id))
	if err != nil {
		log.Printf("Error downloading demo: %v\n", err)
		return err
	}

	demo, err := os.Create(fmt.Sprintf("%s/portal2/demos/autorender/%s.dem", config.Cfg.GameDir, id))
	if err != nil {
		log.Println("error creating demo file", err)
		return err
	}
	defer demo.Close()

	_, err = io.Copy(demo, downloadResponse.Body)
	if err != nil {
		log.Println("error writing demo", err)
		return err
	}

	return nil
}

func PrepareAutoexec(id string) error {
	autoexec, err := os.Create(fmt.Sprintf("%s/portal2/cfg/autoexec.cfg", config.Cfg.GameDir))
	if err != nil {
		log.Printf("Error creating autoexec.cfg: %v\n", err)
		return err
	}

	autoexec.WriteString("plugin_load sar\nexec worker.cfg\n")
	autoexec.WriteString(fmt.Sprintf("playdemo demos/autorender/%s\n", id))
	autoexec.WriteString(fmt.Sprintf("sar_render_start autorender/%s.mp4\n", id))

	autoexec.Close()
	return nil
}

func DownloadConfigs() error {
	response, err := http.Get(fmt.Sprintf("%s/storage/worker.cfg", config.Cfg.Server))
	if err != nil || response.Status == "404 NOT FOUND" {
		log.Printf("Error downloading worker.cfg: %v\n", err)
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s/portal2/cfg/worker.cfg", config.Cfg.GameDir))
	if err != nil {
		log.Println("error creating worker.cfg", err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Println("error writing worker.cfg", err)
		return err
	}

	return nil
}

func UploadVideo(id string) error {
	video, err := os.Open(fmt.Sprintf("%s/portal2/demos/autorender/%s.dem.mp4", config.Cfg.GameDir, id))
	if err != nil {
		log.Printf("error opening video: %v\n", err)
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	formValue, _ := writer.CreateFormFile("file", fmt.Sprintf("%s.mp4", id))
	io.Copy(formValue, video)

	writer.WriteField("videoId", id)
	writer.Close()

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/upload/video", config.Cfg.Server), body)
	if err != nil {
		log.Printf("error creating request: %v\n", err)
		return err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Worker-ID", config.Cfg.Id)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("error uploading video: %v\n", err)
		return err
	}

	if response.Status != "200 OK" {
		log.Printf("server error uploading video")
	}

	return nil
}

func DeleteFiles() error {
	contents, err := filepath.Glob(fmt.Sprintf("%s/demos/autorender", config.Cfg.GameDir))
	if err != nil {
		return err
	}

	for _, file := range contents {
		err = os.RemoveAll(file)
		if err != nil {
			log.Printf("Failed to delete file: %s\n", file)
			continue
		}
	}

	return nil
}
