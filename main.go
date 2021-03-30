package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	urls := make(map[string]string)
	args := os.Args[1:]

	if len(args) > 1 {
		if filepath.IsAbs(strings.TrimSpace(args[1])) {
			urls[strings.TrimSpace(args[0])] = strings.TrimSpace(args[1])
		} else {
			urls[strings.TrimSpace(args[0])] = path.Join(GetExeDirectory(), strings.TrimSpace(args[1]))
		}
	} else if len(args) == 1 {
		u, err := url.Parse(args[0])
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		urls[strings.TrimSpace(args[0])] = path.Join(GetExeDirectory(), strings.TrimSpace(u.Path))
	} else {
		file, err := os.Open("mangas.txt")
		if err != nil {
			log.Fatal(err)
			os.Exit(3)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			s := strings.Split(scanner.Text(), ";")
			urls[strings.TrimSpace(s[0])] = strings.TrimSpace(s[1])
		}
	}

	// urls := map[string]string{
	// 	"https://official-complete-1.granpulse.us/manga/Akira/":      "D:\\Temp\\!Manga\\Akira\\",
	// 	"https://scans-complete.hydaelyn.us/manga/Gantz/":            "D:\\Temp\\!Manga\\Gantz\\",
	// 	"https://official-complete-2.eorzea.us/manga/Toukyou-Kushu/": "D:\\Temp\\!Manga\\Tokyo Ghoul\\",
	// 	"https://official-complete-1.granpulse.us/manga/Uzumaki/":    "D:\\Temp\\!Manga\\Uzumaki\\",
	// }

	var wg sync.WaitGroup

	for k, v := range urls {
		wg.Add(1)
		baseUrl := k
		basePath := v
		chapterStart := 5
		chapterCount := 1

		go DownloadManga(chapterStart, chapterCount, basePath, baseUrl, &wg)
	}
	wg.Wait()
}

func GetExeDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
		os.Exit(10)
	}
	return filepath.Dir(ex)
}

func DownloadManga(chapterStart int, chapterCount int, basePath string, baseUrl string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Manga url: %s\n", baseUrl)
	fmt.Printf("Folder path: %s\n", basePath)

	for chapter := chapterStart; chapter < chapterStart+chapterCount; chapter++ {
		page := 0
		dirPath := path.Join(basePath, fmt.Sprintf("%04d", chapter))

		fmt.Printf("Rozpoczynam pobieranie rodziału %04d\n", chapter)
		err := os.MkdirAll(dirPath, 0777)
		if err != nil {
			fmt.Println("Problem z utworzeniem katalogu ", dirPath, err)
			break
		}

		for {
			page++
			url := fmt.Sprintf("%s%04d-%03d.png", baseUrl, chapter, page)
			filePath := path.Join(basePath, fmt.Sprintf("%04d\\%03d.png", chapter, page))

			fmt.Printf("Pobieram plik %s\n", url)
			err := DownloadFile(filePath, url)
			if err != nil {
				fmt.Println(err)
				break
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
