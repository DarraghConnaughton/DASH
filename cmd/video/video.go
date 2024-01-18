package video

import (
	"dash/cmd/helper"
	"dash/cmd/hlsmanifest"
	"dash/cmd/types"
	"fmt"
	"os"
	"strings"
)

type Video struct {
	VideoUID               string
	EncodedRepresentations []*types.HLSManifest
}

func (v *Video) LoadManifests(resolutions []string) error {
	var manifests []*types.HLSManifest
	paths, err := helper.RecursiveDirectorySearch(fmt.Sprintf("./data/%s", v.VideoUID))
	if err != nil {
		return err
	}
	for _, fp := range paths {
		if helper.Contains(resolutions, fp) {
			tmp := strings.Split(fp, "/")
			quality := tmp[len(tmp)-1]
			file, err := os.Open(fmt.Sprintf("%s/%s.m3u8", fp, quality))
			if err != nil {
				return err
			}

			hlsPlaylist, err := hlsmanifest.ParseHLSManifest(file, fp)
			if err != nil {
				return err
			}
			if err := file.Close(); err != nil {
				return err
			}
			manifests = append(manifests, hlsPlaylist)
		}
	}
	v.EncodedRepresentations = manifests
	return nil
}
