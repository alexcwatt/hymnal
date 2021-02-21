/*
Copyright © 2020 Alex Watt <alex@alexcwatt.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"

	"hymnal/data"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/spf13/cobra"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play a hymn",
	Long: `Play a hymn by number or at random. For example:

hymnal play 100
hymnal play random`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		if args[0] == "random" {
			return nil;
		}
		if _, err := strconv.Atoi(args[0]); err != nil {
			return errors.New("must be an integer or 'random'")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var hymn int
		if args[0] == "random" {
			rand.Seed(time.Now().Unix())
			hymn = data.Hymns[rand.Intn(len(data.Hymns))].Number
		} else {
			hymn, _ = strconv.Atoi(args[0])
		}
		preview, _ := cmd.Flags().GetBool("preview")
		play(hymn, preview)
	},
}

func play(hymn int, preview bool) error {
	archive, err := zip.OpenReader(data.ZipPath)
	if err != nil {
		fmt.Println("Unable to read MP3 zip.")
		return err
	}
	for _, f := range archive.File {
		if f.Name[4:7] == fmt.Sprintf("%03d", hymn) {
			hymn_data := data.Hymns[hymn-1]
			fmt.Printf("Playing Hymn %s\n", hymn_data)
			rc, err := f.Open()
			if err != nil {
				fmt.Println("Unable to open MP3 file.")
				return err
			}
			playFile(rc, preview)
		}
	}
	return nil
}

func playFile(file io.ReadCloser, preview bool) {
	streamer, format, err := mp3.Decode(file)
	if err != nil {
		fmt.Println("Unable to decode MP3 file.")
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	volume := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	speaker.Play(beep.Seq(volume, beep.Callback(func() {
		done <- true
	})))

	if preview {
		go func() {
			time.Sleep(3 * time.Second)
			for i := 0; i < 20; i++ {
				volume.Volume -= 0.1
				time.Sleep(100 * time.Millisecond)
			}
			done <- true
		}()
	}

	<-done
}

func init() {
	rootCmd.AddCommand(playCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	playCmd.Flags().BoolP("preview", "p", false, "Play a 5-second preview")
}
