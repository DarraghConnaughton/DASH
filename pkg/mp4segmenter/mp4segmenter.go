package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func EncodeAndSegmentMP4(inputpath string, outputpath string) {
	resolutions := []string{"1080p", "720p", "480p", "360p"}
	if _, err := os.Stat(outputpath); os.IsNotExist(err) {
		if err := os.MkdirAll(outputpath, 0755); err != nil {
			log.Println(err)
			return
		}

		for _, res := range resolutions {
			if err := os.MkdirAll(fmt.Sprintf("%s/%s", outputpath, res), 0755); err != nil {
				log.Println(err)
				return
			}
		}
	}

	if _, err := os.Stat(inputpath); os.IsNotExist(err) {
		return
	}
	ffmpegCommand := fmt.Sprintf(`
    ffmpeg -i %s -hls_list_size 0 -vf "scale=1920:1080" -c:v h264 -b:v 5000k -c:a aac -b:a 192k %s/1080p/1080p.m3u8 \
           -hls_list_size 0 -vf "scale=1280:720" -c:v h264 -b:v 2500k -c:a aac -b:a 128k %s/720p/720p.m3u8 \
           -hls_list_size 0 -vf "scale=854:480" -c:v h264 -b:v 1200k -c:a aac -b:a 96k %s/480p/480p.m3u8 \
           -hls_list_size 0 -vf "scale=640:360" -c:v h264 -b:v 800k -c:a aac -b:a 64k %s/360p/360p.m3u8
`, inputpath, outputpath, outputpath, outputpath, outputpath)

	cmd := exec.Command("sh", "-c", ffmpegCommand)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Println("Error:", err)
		return
	}
	log.Println("Command output:", string(output))
}

func main() {
	EncodeAndSegmentMP4("video2.mp4", "./data/video2")
}
