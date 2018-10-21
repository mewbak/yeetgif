package main

import (
	"bytes"
	"log"
	"sort"
	"strings"

	"github.com/sgreben/yeetgif/pkg/imaging"

	"github.com/gobwas/glob"
	"github.com/sgreben/yeetgif/pkg/gifcmd"
	"github.com/sgreben/yeetgif/pkg/gifstatic"

	cli "github.com/jawher/mow.cli"
)

func CommandEmoji(cmd *cli.Cmd) {
	cmd.Spec = "[OPTIONS] EMOJI..."
	var (
		size       = gifcmd.FloatsCSV{Values: []float64{128}}
		alpha      = gifcmd.FloatsCSV{Values: []float64{1.0}}
		queryParts = cmd.StringsArg("EMOJI", nil, "one or more glob expressions")
		exact      = cmd.BoolOpt("e exact", false, "match the query exactly")
		pipe       = cmd.BoolOpt("p pipe", false, "overlay the emoji over input images (instead of just creating one)")
		list       = cmd.BoolOpt("l list-only", false, "just list matches")
	)
	cmd.VarOpt("s size", &size, "")
	cmd.VarOpt("a pipe-alpha", &alpha, "")
	cmd.Action = func() {
		if *pipe {
			ProcessInput()
		}
		queryBuffer := &bytes.Buffer{}
		if !*exact {
			queryBuffer.WriteRune('*')
		}
		for i, part := range *queryParts {
			queryBuffer.WriteString(part)
			if i < len(*queryParts)-1 {
				queryBuffer.WriteRune('*')
			}
		}
		if !*exact {
			queryBuffer.WriteRune('*')
		}
		query := queryBuffer.String()
		queryGlob, err := glob.Compile(query)
		if err != nil {
			log.Fatalf("parse glob %q: %v", query, err)
		}
		matches := EmojiMatches(queryGlob)
		if len(matches) == 0 {
			log.Fatalf("no emoji matches found for %q", query)
		}
		if *list || len(matches) > 1 {
			log.Printf("%d matching emoji found for %q", len(matches), query)
			for _, m := range matches {
				log.Printf("%s\t%s", string(m.Runes), m.UnicodeNames)
			}
		}
		if *list {
			config.NoOutput = true
			return
		}
		emoji := matches[0]
		log.Printf("picked %s %s", string(emoji.Runes), emoji.UnicodeNames)
		Emoji(emoji, size.PiecewiseLinear(0, 1), alpha.PiecewiseLinear(0, 1))
	}
}

func EmojiMatches(queryGlob glob.Glob) (matches []gifstatic.Emoji) {
	for _, e := range gifstatic.EmojiList {
		if queryGlob.Match(strings.Join(e.UnicodeNames, " ")) {
			matches = append(matches, e)
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		if len(matches[i].UnicodeNamesJoined) < len(matches[j].UnicodeNamesJoined) {
			return true
		}
		if len(matches[i].UnicodeNamesJoined) > len(matches[j].UnicodeNamesJoined) {
			return false
		}
		return matches[i].UnicodeNamesJoined < matches[j].UnicodeNamesJoined
	})
	return
}

func Emoji(emoji gifstatic.Emoji, sizeF, alphaF func(float64)float64) {
	size0 := sizeF(0)
	emojiImage := imaging.Resize(emoji.Image(), int(size0), int(size0), imaging.Lanczos)
	if len(images) == 0 {
		images = append(images, emojiImage)
		return
	}
	n := float64(len(images))
	write := func(i int) {
		t := float64(i)/n
		size := sizeF(t)
		if size0 != size {
			emojiImage = imaging.Resize(emoji.Image(), int(size), int(size), imaging.Lanczos)
		}
		alpha := alphaF(t)
		images[i] = imaging.OverlayCenter(images[i], emojiImage, alpha)
	}
	parallel(len(images), write)
}
