package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	// Sample text or you could read from a file
	text := `Go is an open source programming language that makes it easy to build 
    simple, reliable, and efficient software. Go was designed at Google in 2007 
    by Robert Griesemer, Rob Pike, and Ken Thompson. Go is syntactically similar 
    to C, but with memory safety, garbage collection, structural typing, and 
    CSP-style concurrency. The language is often referred to as Golang because of 
    its former domain name, golang.org, but the proper name is Go.`

	// Convert to lowercase and split into words
	text = strings.ToLower(text)

	// Use regex to extract words
	re := regexp.MustCompile(`[a-z]+`)
	words := re.FindAllString(text, -1)

	// Count word frequencies
	frequencies := make(map[string]int)
	for _, word := range words {
		frequencies[word]++
	}

	// Filter out common words (simplified stop words list)
	stopWords := map[string]bool{
		"the": true, "and": true, "is": true, "to": true, "of": true,
		"a": true, "in": true, "but": true, "with": true, "by": true,
		"was": true, "its": true,
	}

	// Create a list of word-frequency pairs
	type WordFreq struct {
		Word  string
		Count int
	}

	var wordFreqs []WordFreq
	for word, count := range frequencies {
		if !stopWords[word] && len(word) > 1 {
			wordFreqs = append(wordFreqs, WordFreq{word, count})
		}
	}

	// Sort by frequency (descending)
	sort.Slice(wordFreqs, func(i, j int) bool {
		return wordFreqs[i].Count > wordFreqs[j].Count
	})

	// Print top N words
	topN := 10
	if len(wordFreqs) < topN {
		topN = len(wordFreqs)
	}

	fmt.Printf("Top %d words:\n", topN)
	fmt.Printf("%-15s %s\n", "WORD", "FREQUENCY")
	fmt.Println("------------------------")

	for i := 0; i < topN; i++ {
		wf := wordFreqs[i]
		fmt.Printf("%-15s %d\n", wf.Word, wf.Count)
	}

	// Print to file (optional)
	file, err := os.Create("word_frequencies.txt")
	if err == nil {
		defer file.Close()

		fmt.Fprintf(file, "%-15s %s\n", "WORD", "FREQUENCY")
		fmt.Fprintln(file, "------------------------")

		for _, wf := range wordFreqs {
			fmt.Fprintf(file, "%-15s %d\n", wf.Word, wf.Count)
		}
	}
}
