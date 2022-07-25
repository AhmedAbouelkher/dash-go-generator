package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type resolution struct {
	width    int
	height   int
	crf      int
	vBitRate int
}

func (r *resolution) String() string {
	return fmt.Sprintf("%d:%d", r.width, r.height)
}

func (r *resolution) Name() string {
	return fmt.Sprintf("%dp", r.width)
}

var resolutions = []resolution{
	{240, 426, 22, 400},
	{360, 640, 24, 800},
	{480, 842, 24, 1400},
	{720, 1280, 25, 2800},
	{1080, 1920, 26, 5000},
}

func generateResolutions(ctx context.Context, i string, tf string) ([]string, error) {
	var out []string
	cw, _, err := sourceRes(ctx, i)
	if err != nil {
		return nil, err
	}

	log.Println("Source width:", cw)

	for _, r := range resolutions {
		if r.width > cw {
			continue
		}
		outF, err := generateVideoRes(ctx, r, i, tf)
		if err != nil {
			return nil, err
		}
		log.Printf("Generated %s", outF)
		out = append(out, outF)
	}
	return out, nil
}

func generateResolutionsConcurrent(ctx context.Context, i string, tf string) ([]string, error) {
	// start time
	start := time.Now()

	cw, _, err := sourceRes(ctx, i)
	if err != nil {
		return nil, err
	}

	log.Println("Source width:", cw)

	rls := filter(resolutions, func(r resolution) bool {
		return r.width <= cw
	})

	res := RunWorkersPool(&Pool[resolution, string]{
		Workers: 2,
		Data:    rls,
		Consumer: func(r resolution) string {
			// log generating resolution
			log.Printf("Generating %s", r.String())
			outF, err := generateVideoRes(ctx, r, i, tf)
			if err != nil {
				log.Printf("Error in [generateVideoRes] %s", err)
				return ""
			}
			log.Printf("Generated %s", outF)
			return outF
		},
	})

	// end time
	end := time.Now()
	log.Printf("Time taken: %fs", end.Sub(start).Seconds())
	return res, nil
}

func sourceRes(ctx context.Context, i string) (int, int, error) {
	c := fmt.Sprintf("ffprobe -v error -select_streams v:0 -show_entries stream=width,height -of csv=s=x:p=0 %s", i)
	args := strings.Split(c, " ")

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	log.Println(args)

	o, err := cmd.CombinedOutput()

	if err != nil {
		log.Print("Error in [sourceRes]", string(o))
		return 0, 0, err
	}

	w, h, pE := parseRes(strings.Trim(string(o), "\n"))
	if pE != nil {
		return 0, 0, pE
	}
	return w, h, nil
}

func parseRes(res string) (int, int, error) {
	v := strings.Split(res, "x")
	if len(v) != 2 {
		return 0, 0, fmt.Errorf("invalid resolution %s", res)
	}
	w, err := strconv.ParseInt(v[1], 10, 32)
	if err != nil {
		return 0, 0, err
	}
	h, err := strconv.ParseInt(v[0], 10, 32)
	if err != nil {
		return 0, 0, err
	}
	return int(w), int(h), nil
}

func generateVideoRes(ctx context.Context, r resolution, i string, tf string) (string, error) {
	outF := fmt.Sprintf("%s_%s.mp4", tf, r.Name())

	cmn := `-preset slow -tune film -vsync passthrough -an -c:v libx264`
	c := fmt.Sprintf(
		`ffmpeg -y -i %s %s -crf %d -b:v %dk -pix_fmt yuv420p -vf scale=%s -f mp4 %s`,
		i,
		cmn,
		r.crf,
		r.vBitRate,
		r.String(),
		outF,
	)

	args := strings.Split(c, " ")

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	log.Println(args)

	o, err := cmd.CombinedOutput()

	if err != nil {
		log.Print(string(o))
		return "", err
	}
	return outF, nil
}

func filter[T any](d []T, cond func(v T) bool) []T {
	var out []T
	for _, v := range d {
		if cond(v) {
			out = append(out, v)
		}
	}
	return out
}
