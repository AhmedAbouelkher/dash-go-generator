package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

/*
MP4Box -dash 2000 -rap -frag-rap  -bs-switching no -profile "dashavc264:live" "${targetFile}_240p.mp4" "${targetFile}_360p.mp4" "${targetFile}_480p.mp4" "${targetFile}_720p.mp4" "${targetFile}_1080p.mp4" "${targetFile}_audio.m4a" -out "${target}/playlist.mpd"
rm "${targetFile}_240p.mp4" "${targetFile}_360p.mp4" "${targetFile}_480p.mp4" "${targetFile}_720p.mp4" "${targetFile}_1080p.mp4" "${targetFile}_audio.m4a"
*/

func generateDash(ctx context.Context, res []string, a string, t string) error {
	b := `MP4Box -dash 10000 -frag 1000 -rap -frag-rap -bs-switching no -profile dashavc264:live`
	i := strings.Join(append(res, a), " ")
	out := fmt.Sprintf("%s/playlist.mpd", t)
	c := fmt.Sprintf("%s %s -out %s", b, i, out)

	args := strings.Split(c, " ")

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	log.Println(cmd.Args)

	o, err := cmd.CombinedOutput()

	if err != nil {
		log.Print(string(o))
		return err
	}

	if err := cleanFiles(ctx, res, a); err != nil {
		return err
	}

	return nil
}

func cleanFiles(ctx context.Context, res []string, a string) error {
	i := strings.Join(append(res, a), " ")
	c := fmt.Sprintf("rm %s", i)

	args := strings.Split(c, " ")

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	log.Println(cmd.Args)

	o, err := cmd.CombinedOutput()

	if err != nil {
		log.Print(string(o))
		return err
	}
	return nil
}
