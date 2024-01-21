package hlsmanifest

import (
	"bufio"
	"dash/pkg/types"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseHLSManifest(file *os.File, filepath string) (*types.HLSManifest, error) {
	scanner := bufio.NewScanner(file)
	hlsPlaylist := &types.HLSManifest{}
	segmentCount := 0
	segmentMap := make(map[int]types.Segment)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#EXT-X-VERSION:") {
			version, err := strconv.Atoi(strings.TrimPrefix(line, "#EXT-X-VERSION:"))
			if err != nil {
				return nil, err
			}
			hlsPlaylist.Version = version
		} else if strings.HasPrefix(line, "#EXT-X-TARGETDURATION:") {
			targetDuration, err := strconv.Atoi(strings.TrimPrefix(line, "#EXT-X-TARGETDURATION:"))
			if err != nil {
				return nil, err
			}
			hlsPlaylist.TargetDuration = targetDuration
		} else if strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:") {
			mediaSequence, err := strconv.Atoi(strings.TrimPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"))
			if err != nil {
				return nil, err
			}
			hlsPlaylist.MediaSequence = mediaSequence
		} else if strings.HasPrefix(line, "#EXTINF:") {
			// Extract duration from the line
			duration := strings.TrimSuffix(strings.TrimPrefix(line, "#EXTINF:"), ",")

			// Create a new Segment instance
			segmentMap[segmentCount] = types.Segment{
				Location: fmt.Sprintf("%s/%s", filepath, scanner.Scan()),
				Duration: duration,
			}
			segmentCount += 1

		}
		//else if strings.Contains(line, ".ts") {
		//	segmentMap[segmentCount] = fmt.Sprintf("%s/%s", filepath, line)
		//	segmentCount += 1
		//}
	}
	hlsPlaylist.SegmentLocations = segmentMap
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return hlsPlaylist, nil
}
