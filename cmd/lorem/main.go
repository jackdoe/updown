package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"math/rand"
	"os"
)

var wordsList = [][]byte{
	[]byte("ipsum"), []byte("semper"), []byte("habeo"), []byte("duo"), []byte("ut"), []byte("vis"), []byte("aliquyam"), []byte("eu"), []byte("splendide"), []byte("Ut"), []byte("mei"), []byte("eteu"), []byte("nec"), []byte("antiopam"), []byte("corpora"), []byte("kasd"), []byte("pretium"), []byte("cetero"), []byte("qui"), []byte("arcu"), []byte("assentior"), []byte("ei"), []byte("his"), []byte("usu"), []byte("invidunt"), []byte("kasd"), []byte("justo"), []byte("ne"), []byte("eleifend"), []byte("per"), []byte("ut"), []byte("eam"), []byte("graeci"), []byte("tincidunt"), []byte("impedit"), []byte("temporibus"), []byte("duo"), []byte("et"), []byte("facilisis"), []byte("insolens"), []byte("consequat"), []byte("cursus"), []byte("partiendo"), []byte("ullamcorper"), []byte("Vulputate"), []byte("facilisi"), []byte("donec"), []byte("aliquam"), []byte("labore"), []byte("inimicus"), []byte("voluptua"), []byte("penatibus"), []byte("sea"), []byte("vel"), []byte("amet"), []byte("his"), []byte("ius"), []byte("audire"), []byte("in"), []byte("mea"), []byte("repudiandae"), []byte("nullam"), []byte("sed"), []byte("assentior"), []byte("takimata"), []byte("eos"), []byte("at"), []byte("odio"), []byte("consequat"), []byte("iusto"), []byte("imperdiet"), []byte("dicunt"), []byte("abhorreant"), []byte("adipisci"), []byte("officiis"), []byte("rhoncus"), []byte("leo"), []byte("dicta"), []byte("vitae"), []byte("clita"), []byte("elementum"), []byte("mauris"), []byte("definiebas"), []byte("uonsetetur"), []byte("te"), []byte("inimicus"), []byte("nec"), []byte("mus"), []byte("usu"), []byte("duo"), []byte("aenean"), []byte("corrumpit"), []byte("aliquyam"), []byte("est"), []byte("eum"),
}

var space = []byte(" ")
var lorem = []byte("Lorem ")

func getRandomWord() []byte {
	return wordsList[rand.Intn(len(wordsList))]
}

func generateWords(fd io.Writer, length int) error {
	_, err := fd.Write(lorem)
	if err != nil {
		return err
	}

	for i := 0; i < length-1; i++ {
		_, err := fd.Write(getRandomWord())
		if err != nil {
			return err
		}
		if i != length-2 {
			_, err = fd.Write(space)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func generateParagraphs(fd io.Writer, count, length int, separator []byte) error {
	if length == 0 {
		length = 10
	}
	for i := 0; i < count; i++ {
		err := generateWords(fd, length)
		if err != nil {
			return err
		}
		_, err = fd.Write(separator)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var paragraphs = flag.Int("p", 0, "how many paragraphs to generate")
	var words = flag.Int("w", 0, "how many words to generate")
	var separator = flag.String("s", "\n", "the separator between paragraphs")

	flag.Parse()

	writer := bufio.NewWriter(os.Stdout)

	if *paragraphs == 0 && *words == 0 {
		flag.Usage()
	} else if *paragraphs != 0 {
		err := generateParagraphs(writer, *paragraphs, *words, []byte(*separator))
		if err != nil {
			log.Fatal("error writing, err: %s", err.Error())
		}
	} else if *words != 0 {
		err := generateWords(writer, *words)
		if err != nil {
			log.Fatal("error writing, err: %s", err.Error())
		}
	}
	writer.Flush()
}
