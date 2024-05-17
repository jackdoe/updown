package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"9fans.net/go/acme"
	"path"
	"strconv"
)

func ListSlides() []string {
	fn := os.Getenv("%")
	if fn == "" {
		panic("need acme")
	}

	dir := filepath.Dir(fn)
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	slides := []string{}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".txt") {
			slides = append(slides, path.Join(dir, f.Name()))
		}
	}
	sort.Strings(slides)

	return slides
}

const SLIDE_PLUS = 1
const SLIDE_MINUS = 2

func main() {
	wid, _ := strconv.ParseInt(os.Getenv("winid"), 10, 32)
	w, err := acme.Open(int(wid), nil)
	if err != nil {
		panic(err)
	}
	direction := SLIDE_PLUS
	if strings.HasSuffix(os.Args[0], "-") {
		direction = SLIDE_MINUS
	}
	slides := ListSlides()
	currentSlide := os.Getenv("%")
	replaceSlide := ""
	i := sort.Search(len(slides), func(i int) bool { return slides[i] >= currentSlide })
	if i < len(slides) && slides[i] == currentSlide {
		if direction == SLIDE_PLUS {
			if i < len(slides)-1 {
				replaceSlide = slides[i+1]
			}
		} else {
			if i > 0 {
				replaceSlide = slides[i-1]
			}
		}
	}
	if replaceSlide != "" && replaceSlide != currentSlide {
		b, err := os.ReadFile(replaceSlide)
		if err != nil {
			panic(err)
		}
		w.Name(replaceSlide)
		w.Clear()
		w.Addr(",")
		w.Write("data", b)
	}
}
