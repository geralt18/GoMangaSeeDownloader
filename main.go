package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("mangas.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}
	defer file.Close()

	urls := make(map[string]string)
	// urls := map[string]string{
	// 	//"https://official-complete-1.granpulse.us/manga/Akira/":      "D:\\Temp\\!Manga\\Akira\\",
	// 	"https://scans-complete.hydaelyn.us/manga/Gantz/":            "D:\\Temp\\!Manga\\Gantz\\",
	// 	"https://official-complete-2.eorzea.us/manga/Toukyou-Kushu/": "D:\\Temp\\!Manga\\Tokyo Ghoul\\",
	// 	"https://official-complete-1.granpulse.us/manga/Uzumaki/":    "D:\\Temp\\!Manga\\Uzumaki\\",
	// }

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ";")
		urls[strings.TrimSpace(s[0])] = strings.TrimSpace(s[1])
	}

	for k, v := range urls {
		baseUrl := k
		basePath := v
		chapterStart := 5
		chapterCount := 1

		for chapter := chapterStart; chapter < chapterStart+chapterCount; chapter++ {
			page := 0
			dirPath := fmt.Sprintf("%s%04d", basePath, chapter)

			fmt.Printf("Rozpoczynam pobieranie rodziału %04d\n", chapter)
			err := os.MkdirAll(dirPath, 0777)
			if err != nil {
				fmt.Println("Problem z utworzeniem katalogu ", dirPath, err)
				break
			}

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
