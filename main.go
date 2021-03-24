package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	urls := map[string]string{
		//"https://official-complete-1.granpulse.us/manga/Akira/":      "D:\\Temp\\!Manga\\Akira\\",
		"https://scans-complete.hydaelyn.us/manga/Gantz/":            "D:\\Temp\\!Manga\\Gantz\\",
		"https://official-complete-2.eorzea.us/manga/Toukyou-Kushu/": "D:\\Temp\\!Manga\\Tokyo Ghoul\\",
		"https://official-complete-1.granpulse.us/manga/Uzumaki/":    "D:\\Temp\\!Manga\\Uzumaki\\",
	}
	for k, v := range urls {
		baseUrl := k
		basePath := v
		chapterCount := 5

		for chapter := 1; chapter < chapterCount; chapter++ {
			fmt.Printf("Rozpoczynam pobieranie rodziału %04d\n", chapter)
			os.MkdirAll(fmt.Sprintf("%s%04d", basePath, chapter), 0777)
			page := 0

			for {
				page++
				url := fmt.Sprintf("%s%04d-%03d.png", baseUrl, chapter, page)
				filePath := fmt.Sprintf("%s%04d\\%03d.png", basePath, chapter, page)

				fmt.Printf("Pobieram plik %s\n", url)
				err := DownloadFile(filePath, url)
				if err != nil {
					fmt.Println(err)
					break
				}
			}
		}
	}
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
	//File existsl
	if _, err := os.Stat(filepath); err == nil {
		return nil
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Błąd pobrania pliku %s. StatusCode=%d", url, resp.StatusCode)
		//return errors.New(fmt.Sprintf("Błąd pobrania pliku %s. StatusCode=%d", url, resp.StatusCode))
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
