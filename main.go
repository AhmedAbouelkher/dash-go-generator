package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	ctx = context.Background()
)

func main() {
	var i, t, typ string

	flag.StringVar(&i, "i", "", "input file")
	flag.StringVar(&t, "t", ".", "target")
	flag.StringVar(&typ, "type", "dash", "type of output [dash, hls or both]")

	flag.Parse()
	_, err := generate(ctx, i, t, typ)
	if err != nil {
		panic(err)
	}
}

func generate(ctx context.Context, i string, t string, typ string) (string, error) {
	// TODO: Check if input file exist
	log.Println(t, i)
	if err := os.RemoveAll(t); err != nil {
		return "", err
	}

	// Validate target [t] directory and create it
	if err := os.MkdirAll(t, 0777); err != nil {
		return "", err
	}

	tf := targetFile(i, t)

	log.Println("target file", tf)

	// Get the input file audio bit rate
	abr, abrErr := audioBitRate(ctx, i)
	if abrErr != nil {
		return "", abrErr
	}
	log.Println(abr)

	// Generate [audio] using ffmpeg
	af, aErr := generateAudio(ctx, i, tf, abr)
	if aErr != nil {
		return "", aErr
	}
	log.Println(af)

	rsl, resErr := generateResolutionsConcurrent(ctx, i, tf)
	if resErr != nil {
		return "", resErr
	}

	log.Println("Finished", rsl)

	if len(rsl) == 0 {
		return "", nil
	}

	// Generate different resolutions.
	// Comping resolutions using mp4box
	// delete all temp *.mp4 files
	// return the output path of playlist.mpd
	switch typ {
	case "dash":
		if err := generateDash(ctx, rsl, af, t); err != nil {
			return "", err
		}
		log.Println("Finished Dash")

	case "hls":
		if err := generateHls(ctx, rsl, af, t); err != nil {
			return "", err
		}
		log.Println("Finished HLS")
	case "both":
		if err := generateHlsAndDash(ctx, rsl, af, t); err != nil {
			return "", err
		}

		log.Println("Finished hls and dash")

	default:
		return "", fmt.Errorf("unknown type %s", typ)

	}

	return "", nil
}

func targetFile(i string, t string) string {
	f := filepath.Base(filepath.Dir(i))
	return filepath.Join(t, f)
}
