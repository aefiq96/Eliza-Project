//Nuzhafiq Iqbal - G00324934
//adapted from https://github.com/data-representation/eliza
package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Replacer is a struct with two elements
type Replacer struct {
	original     *regexp.Regexp
	replacing []string
}

// ReadReplacersFromFile reads an array of Replacers from a text file.
func ReadReplacersFromFile(path string) []Replacer {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var replacers []Replacer

	// Read the file line by line.
	for scanner, readoriginal := bufio.NewScanner(file), false; scanner.Scan(); {
		switch line := scanner.Text(); {
		case strings.HasPrefix(line, "#"):
		case len(line) == 0:
			readoriginal = false		
		case readoriginal == false:
			replacers = append(replacers, Replacer{original: regexp.MustCompile(line)})
			readoriginal = true
		default:
			replacers[len(replacers)-1].replacing = append(replacers[len(replacers)-1].replacing, line)
		}
	}
	// Return the replacers array.
	return replacers
}
// Eliza is a chatbot.
// The fields answers and subs arrays of Replacers.
type Eliza struct {
	answers     []Replacer
	subs []Replacer
}

func ElizaFromFiles(responsePath string, substitutionPath string) Eliza {
	eliza := Eliza{}

	eliza.answers = ReadReplacersFromFile(responsePath)
	eliza.subs = ReadReplacersFromFile(substitutionPath)

	return eliza
}
// Responds to take a string as input and returns a string. The returned string
// contains the chatbot's response to the input.
func (me *Eliza) RespondTo(input string) string {
	for _, respond := range me.answers {
		if matches := respond.original.FindStringSubmatch(input); matches != nil {
			output := respond.replacing[rand.Intn(len(respond.replacing))]
			boundary := regexp.MustCompile(`[\s,.?!]+`)
			for m, match := range matches[1:] {
				tokens := boundary.Split(match, -1)
				for t, token := range tokens {
					for _, subway := range me.subs {
						if subway.original.MatchString(token) {
							tokens[t] = subway.replacing[rand.Intn(len(subway.replacing))]
							break
						}
					}
				}
				output = strings.Replace(output, "$"+strconv.Itoa(m+1), strings.Join(tokens, " "), -1)
			}
			// Send the filled answers back.
			return output
		}
	}
	// This returns if there are no matches
	return "I don't know what to say."
}

func userinputhandler(w http.ResponseWriter, r *http.Request) {

	eliza := ElizaFromFiles("data/answers.txt", "data/subs.txt")
	// Read from the user.
	input := r.URL.Query().Get("value")
	//print the output
	fmt.Fprintf(w, "User: %s\n", input)
	fmt.Fprintf(w, "Eliza: %s\n", eliza.RespondTo(input))

}
// Program entry point.
func main() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/user-input", userinputhandler)
	http.ListenAndServe(":8080", nil)
}
