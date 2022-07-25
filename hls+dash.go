package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// gpac -i data_240p.mp4 -i data_360p.mp4 -i data_audio.m4a -o res/live.mpd:dual:profile=live
// MP4Box -dash 10000 -frag 1000 -rap -frag-rap -bs-switching no -profile dashavc264:live

func generateHlsAndDash(ctx context.Context, res []string, a string, t string) error {
	b := `MP4Box -dash 10000 -frag 1000 -rap -frag-rap -bs-switching no -profile dashavc264:live`
	i := strings.Join(append(res, a), " ")
	out := fmt.Sprintf("http://ORIG_SERVER_IP_PORT/%s/live.mpd:gpac:dual:llhls=br:rdirs=outdir", t)
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
