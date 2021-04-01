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

var (
	Version   = "0.0.0"
	BuildTime = "2000-01-01T00:00:00"
)

func main() {
	fmt.Println("Version:\t", Version)
	fmt.Println("BuildTime:\t", BuildTime)

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
			log.Println(err)
			fmt.Scanln() // wait for Enter Key
		}
		urls[strings.TrimSpace(args[0])] = path.Join(GetExeDirectory(), strings.TrimSpace(u.Path))
	} else if file, err := os.Open("mangas.txt"); err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			s := strings.Split(scanner.Text(), ";")
			urls[strings.TrimSpace(s[0])] = strings.TrimSpace(s[1])
		}
	} else {
		if err != nil {
			log.Println(err)
			fmt.Scanln() // wait for Enter Key
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
		chapterStart := 1
		chapterCount := 9999

		go DownloadManga(chapterStart, chapterCount, basePath, baseUrl, &wg)
	}
	wg.Wait()
}

func GetExeDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(ex)
}

func DirIsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func DownloadManga(chapterStart int, chapterCount int, basePath string, baseUrl string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Manga url: %s\n", baseUrl)
	fmt.Printf("Folder path: %s\n", basePath)
	chapterErrorCount := 0

	for chapter := chapterStart; chapter < chapterStart+chapterCount; chapter++ {
		page := 0
		chapterPath := path.Join(basePath, fmt.Sprintf("%04d", chapter))

		fmt.Printf("Downloading chapter %04d\n", chapter)
		err := os.MkdirAll(chapterPath, 0777)
		if err != nil {
			fmt.Println("Error creating directory ", chapterPath, err)
			break
		}

		for {
			page++
			url := fmt.Sprintf("%s%04d-%03d.png", baseUrl, chapter, page)
			filePath := path.Join(basePath, fmt.Sprintf("%04d\\%03d.png", chapter, page))

			fmt.Printf("Downloading file %s\n", url)
			err := DownloadFile(filePath, url)
			if err != nil {
				fmt.Println(err)
				if page == 1 {
					chapterErrorCount++
					if v, _ := DirIsEmpty(chapterPath); v {
						os.Remove(chapterPath)
					}
				}

				break
			}
		}

		if chapterErrorCount >= 3 {
			break
		}
	}
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
	//File exists
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
		return fmt.Errorf("error downloading file %s. statuscode=%d", url, resp.StatusCode)
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
