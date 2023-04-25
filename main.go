package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"os"
	"strings"

	"index/suffixarray"
)

// Play represents a play and its quote.
type Play struct {
	Title        string
	Player       string
	ActSceneLine string
	Quote        string
}

// Searcher is used to search for quotes in the Complete Works of Shakespeare.
type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
	Plays         []Play
}

func main() {
	// Load searchers from files
	searcherCSV := Searcher{}
	err := searcherCSV.Load("completeworkssorted.csv")
	if err != nil {
		log.Fatal(err)
	}

	searcherTXT := Searcher{}
	err = searcherTXT.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Setup HTTP server
	http.HandleFunc("/search", handleSearch(searcherCSV))
	http.HandleFunc("/search-context", handleSearchContext(searcherTXT))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// handleSearch is a function that takes a Searcher and returns an http.HandlerFunc.
// It handles incoming HTTP requests and responds with search results.
func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the query parameter from the URL.
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			// If the query parameter is missing, return a bad request error.
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		
		// Perform the search using the provided Searcher.
		results := searcher.Search(query[0])
		
		// Encode the search results as JSON and write the response.
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			// If there is an encoding error, return an internal server error.
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

// handleSearchContext returns a function that handles a search request.
// It takes a Searcher as input, which is used to perform the search.
// The returned function takes an http.ResponseWriter and http.Request,
// extracts the search query from the URL parameters, performs the search,
// and writes the results to the response in JSON format.
func handleSearchContext(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Extract search query from URL parameters
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}

		// Perform search using the provided Searcher
		results := searcher.Search(query[0])

		// Encode search results in JSON format
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}

		// Write encoded search results to response
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

// Load reads data from a file with the specified filename and loads it into the Searcher.
// If the file has a .csv extension, it is assumed to be a file containing play data,
// which is parsed and stored in the Searcher's Plays field.
// If the file has a .txt extension, it is assumed to be a file containing the complete works
// of Shakespeare, which is loaded into the Searcher's CompleteWorks field and used to create
// a SuffixArray for substring searches.
func (s *Searcher) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	defer file.Close()

	if strings.HasSuffix(filename, ".csv") {
		// If file is a CSV, parse play data and store in Searcher's Plays field
		reader := csv.NewReader(file)
		reader.FieldsPerRecord = -1
		records, err := reader.ReadAll()
		if err != nil {
			return fmt.Errorf("Load: %w", err)
		}

		for _, record := range records {
			play := Play{
				Title:        record[1],
				ActSceneLine: record[2],
				Player:       record[3],
				Quote:        record[4],
			}
			s.Plays = append(s.Plays, play)
		}
	} else if strings.HasSuffix(filename, ".txt") {
		// If file is a TXT, load complete works of Shakespeare and create SuffixArray for substring searches
		dat, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("Load: %w", err)
		}
		s.CompleteWorks = string(dat)
		s.SuffixArray = suffixarray.New(dat)
	} else {
		return fmt.Errorf("Load: unsupported file extension")
	}

	return nil
}

// Search searches for a given query string in the complete works and plays loaded into the Searcher.
// It returns a slice of maps, where each map contains information about a single search result.
// If the Searcher has a SuffixArray, it returns the 500 characters surrounding the first occurrence of the query string in the complete works as the "Context" field in the map.
// If the Searcher has a Plays slice, it searches through each play's quotes and returns any quotes that contain the query string, along with the title, player, and act/scene/line information for that quote.
// The function returns an empty slice if there are no matches.
func (s *Searcher) Search(query string) []map[string]string {
	// Create an empty slice of maps to hold the search results
	results := []map[string]string{}

	// If the Searcher has a SuffixArray, use it to perform a substring search on the CompleteWorks string
	if s.SuffixArray != nil {
		index := s.SuffixArray.Lookup([]byte(query), -1)

		// If a match is found, create a map with a context field that contains the 500 characters surrounding the first match
		if index != nil {
			startIndex := index[0] - 125
			if startIndex < 0 {
				startIndex = 0
			}
			endIndex := index[0] + 125
			if endIndex > len(s.CompleteWorks) {
				endIndex = len(s.CompleteWorks)
			}
			result := map[string]string{
				"Context": s.CompleteWorks[startIndex:endIndex],
			}
			results = append(results, result)
		}
	}

	// If the Searcher has Plays, search for the query string in each play's Quote field
	if s.Plays != nil {
		for _, play := range s.Plays {
			words := strings.Fields(play.Quote)
			for _, word := range words {
				if strings.EqualFold(word, query) {
					// If a match is found, create a map with the play's Title, Player, Quote, and ActSceneLine fields
					result := map[string]string{
						"Title":        play.Title,
						"Player":       play.Player,
						"Quote":        play.Quote,
						"ActSceneLine": play.ActSceneLine,
					}
					results = append(results, result)
				}
			}
		}
	}

	// Return the slice of maps containing the search results, or an empty slice if no matches were found
	return results
}
