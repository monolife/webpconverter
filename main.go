package main

import (
	"flag"
	"io"
	"os"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/image/webp"
	"image/jpeg"
)

// convertJPEGToPNG converts from JPEG to PNG.
func convertWebpToJpeg(r io.Reader, w io.Writer) error {
	img, err := webp.Decode(r)
	if err != nil {
		return err
	}
	return jpeg.Encode(w, img, nil)
}

func getPath(path string)(base string){
	ext := filepath.Ext(path)
	base = strings.TrimSuffix(path, ext)
	return
}

func convert(path string)(string){
	baseName := getPath(path)

	webp, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	jpg, err := os.Create(baseName+".jpg")
	if err != nil {
		log.Fatal(err)
	}

	convertWebpToJpeg(webp, jpg)
	return baseName+".jpg"
}

func gen(paths []string) <-chan string {
	out := make(chan string)
	go func() {
		for _, p := range paths {
			out <- p
		}
		close(out)
	}()
	return out
}

func conv(in <-chan string, wg *sync.WaitGroup) {
	go func() {
		wg.Add(1)
		for p := range in {
			convert(p)
		}
		wg.Done()
	}()
}

func main(){
	start := time.Now()
	dir := flag.String("d", "", "directory to process")
	flag.Parse()

	if *dir != ""{
		numCpus := runtime.NumCPU()
		log.Printf("Running on %d CPUs\n", numCpus)
		var wg sync.WaitGroup
		webpFiles,_ := filepath.Glob(*dir+"/*.webp")
		in := gen(webpFiles)
		for i := 0; i < numCpus; i++ {
			conv(in, &wg)
		}
		wg.Wait()
	}else if len(flag.Args()) > 0 {
		webpFilename := flag.Args()[0]
		convert(webpFilename)
	}else{
		log.Print("Nothing selected to convert")
	}

	// Code to measure
	duration := time.Since(start)
	log.Print(duration)

	return
}