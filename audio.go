package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func audioBitRate(ctx context.Context, i string) (int, error) {
	c := "ffprobe -v error -select_streams a:0 -show_entries stream=bit_rate -of csv=s=x:p=0"
	args := strings.Split(c, " ")

	cmd := exec.CommandContext(
		ctx,
		args[0],
		append(args[1:], i)...,
	)

	o, err := cmd.CombinedOutput()

	if err != nil {
		log.Print(string(o))
		return 0, err
	}
	oVal, cErr := strconv.Atoi(strings.Trim(string(o), "\n"))
	if cErr != nil {
		return 0, cErr
	}
	fmt := oVal / 1000
	return int(fmt), nil
}

func generateAudio(ctx context.Context, i string, tf string, ba int) (string, error) {
	c := fmt.Sprintf(
		"ffmpeg -y -i %s -c:a aac -b:a %dk -vn %s_audio.m4a",
		i,
		ba,
		tf,
	)
	args := strings.Split(c, " ")

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	log.Printf("Generating audio: %s %d", tf, ba)

	o, err := cmd.CombinedOutput()

	if err != nil {
		log.Print(string(o))
		return "", err
	}
	return fmt.Sprintf("%s_audio.m4a", tf), nil
}
