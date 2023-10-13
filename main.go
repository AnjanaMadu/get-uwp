package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/cheggaaa/pb/v3"
	"github.com/manifoldco/promptui"
)

func main() {
	prompt := promptui.Prompt{Label: "App Name"}
	result, _ := prompt.Run()

	results, err := searchStore(result)
	if err != nil {
		fmt.Println(err)
		return
	}

	items := make([]string, len(results))
	for i, result := range results {
		items[i] = fmt.Sprintf("%s - %s", result.Title, result.PublisherName)
	}
	prompt2 := promptui.Select{Label: "Select an App", Items: items}
	index, _, _ := prompt2.Run()

	prodId := results[index].ProductId

	files, err := getFiles(prodId)
	if err != nil {
		fmt.Println(err)
		return
	}

	last := files[len(files)-1]

	// Download if the file doesn't exist
	fileName := path.Join("downloads", last.Name)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		err := downloadFile(last.URI, fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Open the file
	exec.Command("pwsh", "-Command", ".\\downloads\\"+last.Name).Start()
}

func downloadFile(url string, filepath string) error {
	// Make the downloads directory
	os.Mkdir("downloads", os.ModePerm)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Get the content length
	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	// Create the progress bar
	bar := pb.New(size)
	bar.Start()

	// Create a proxy reader
	reader := bar.NewProxyReader(resp.Body)

	// Write the body to file
	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}

	// Finish the progress bar
	bar.Finish()

	return nil
}
