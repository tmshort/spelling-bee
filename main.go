package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	var letters string
	var center string
	var wordfile string
	var pangram bool

	flag.StringVar(&letters, "letters", "", "List of non-center letters")
	flag.StringVar(&center, "center", "", "The center letter")
	flag.StringVar(&wordfile, "word-file", "/usr/share/dict/words", "Word list")
	flag.BoolVar(&pangram, "pangram", false, "Only show pangrams")

	flag.Parse()

	if len(letters) != 6 {
		fmt.Println("must provide 6 letters")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if len(center) != 1 {
		fmt.Println("must provide single center letter")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if strings.Contains(letters, center) {
		fmt.Println("center letter must not be in letter list")
		flag.PrintDefaults()
		os.Exit(0)
	}
	counts := make(map[rune]int)
	for _, char := range letters {
		if counts[char] == 1 {
			fmt.Printf("letters must not be duplicated")
			flag.PrintDefaults()
			os.Exit(0)
		}
		counts[char]++
	}

	regstr := fmt.Sprintf("^[%s%s]{4,20}$", center, letters)
	regex := regexp.MustCompile(regstr)

	fmt.Printf("regex: %s\n", regstr)

	file, err := os.Open(wordfile)
	if err != nil {
		fmt.Printf("unable to open word file: %w", err)
		os.Exit(0)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var words []string

	for scanner.Scan() {
		text := scanner.Text()
		if !strings.Contains(text, center) {
			continue
		}
		if regex.MatchString(text) {
			words = append(words, text)
		}
	}
	fmt.Printf("found %d words\n", len(words))
	for _, word := range words {
		if len(word) >= 7 {
			counts := make(map[rune]int)
			for _, c := range word {
				counts[c]++
			}
			if len(counts) == 7 {
				fmt.Println("PANGRAM:", word)
				continue
			}
		}
		if !pangram {
			fmt.Println(word)
		}
	}
}
