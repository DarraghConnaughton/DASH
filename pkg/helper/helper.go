package helper

import (
	"dash/pkg/types"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Contains(files []string, file string) bool {
	for _, f := range files {
		if strings.Contains(file, f) {
			return true
		}
	}
	return false
}

func ListDirectory(path string) ([]string, error) {
	var directories []string

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, filepath.Join(path, entry.Name()))
		}
	}

	return directories, nil
}

func RecursiveDirectorySearch(rootPath string) ([]string, error) {
	var directories []string
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			directories = append(directories, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return directories, nil
}

func Exists(fp string) bool {
	if _, err := os.Stat(fp); err != nil {
		return false
	}
	return true
}

func FormatMetric(hb *types.RPCHeartBeat) string {
	return fmt.Sprintf(`[{
		"metric": "dash_service_monitor",
		"timestamp": %d,
		"value": %d,
		"tags": {
			"service": "%s",
			"type": "active_goroutines"
		}
	}]`, time.Now().Unix(), hb.NumberOfGoroutines, hb.UID)
}

func MonitorErrorChannel(errChan chan error, hardfail bool) error {
	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Println("[-]received an error from the goroutine:", err)
				if hardfail {
					log.Println("[-] hard fail mode enabled, exiting main goroutine.")
					return err
				}
			}
		}
	}
}

func WriteCSV(videoid string, data []types.NetworkTraceData) error {
	file, err := os.Create(fmt.Sprintf("%s.%d.output.csv", videoid, time.Now().UnixNano()))
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	header := []string{"Timestamp", "Bytes", "Sequence", "Resolution"}
	if err := writer.Write(header); err != nil {
		panic(err)
	}

	for _, traceData := range data {
		row := []string{traceData.Timestamp, traceData.Bytes, traceData.Sequence, traceData.Resolution}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	return nil
}
