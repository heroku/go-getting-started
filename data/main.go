package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// AccessToken is a sample access token
// which in real life should be stored and
// read from a data store.
var AccessToken = strconv.Itoa(rand.Int())

func init() {
	_, filename, _, _ := runtime.Caller(0)
	currentPath := path.Dir(filename)
	fullpath := path.Join(currentPath, "./../data", "winners.json")
	LoadFromJSON(fullpath)
}

var winners Winners

// Winners struct which contains
// the list of all winners
type Winners struct {
	Winners []Winner `json:"winners"`
}

// Winner struct which contains a team name
// and corresponding year
type Winner struct {
	Country string `json:"country"`
	Year    int    `json:"year"`
}

// In order to be valid, winner country
// cannot be empty and year cannot be in the past
func (w Winner) isValidWinner() bool {
	currentYear := time.Now().Year()
	isValid := len(w.Country) > 0 && w.Year >= currentYear
	return isValid
}

// LoadFromJSON loads/resets the winners from
// the JSON file
func LoadFromJSON(fullpath string) {
	// Load data from JSON into memory
	jsonFile, err := os.Open(fullpath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &winners)
	// Finish loading JSON
}

// ListAllJSON returns all winners
func ListAllJSON() ([]byte, error) {
	json, err := json.Marshal(winners)
	if err != nil {
		return nil, errors.New("Error marshalling JSON")
	}
	return json, nil
}

// ListAllByYear returns winners by year
func ListAllByYear(yearStr string) ([]byte, error) {
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return nil, errors.New("Cannot convert year to int")
	}
	winnersByYear := Winners{}
	for _, winner := range winners.Winners {
		if winner.Year == year {
			winnersByYear.Winners = []Winner{winner}
			break
		}
	}
	json, err := json.Marshal(winnersByYear)
	if err != nil {
		return nil, errors.New("Error marshalling JSON")
	}

	return json, nil
}

// IsAccessTokenValid implements logic
// for validating access token
func IsAccessTokenValid(token string) bool {
	return token == AccessToken
}

// AddNewWinner adds a new _VALID_ winner to
// the list of winners
func AddNewWinner(payload io.Reader) error {
	if payload == nil {
		return errors.New("Invalid payload")
	}
	var newWinner Winner
	dec := json.NewDecoder(payload)
	err := dec.Decode(&newWinner)

	if err != nil || !newWinner.isValidWinner() {
		return errors.New("Did not add new winner")
	}

	winners.Winners = append(winners.Winners, newWinner)
	return nil
}

// PrintUsage prints test commands to the console
func PrintUsage() {
	usage := `
	
	GETting:

	curl -i http://localhost:8000/
	curl -i http://localhost:8000/winners
	curl -i http://localhost:8000/winners?year=1970
	curl -i http://localhost:8000/winners?year=banana
	
	POSTing with NO access token:
	
	curl -i -X POST %NEXTLINE%
	-d "{\"country\":\"Croatia\", \"year\": 2030}" http://localhost:8000/winners
	
	POSTing with valid access token:
	
	curl -i -X POST %NEXTLINE%
	-H "X-ACCESS-TOKEN: %TOKEN%" %NEXTLINE%
	-d "{\"country\":\"Croatia\", \"year\": 2030}" http://localhost:8000/winners
	
	Then check for the newly added winner
	
	curl -i http://localhost:8000/winners
	
	POSTing with invalid data:

	curl -i -X POST %NEXTLINE%
	-H "X-ACCESS-TOKEN: %TOKEN%" %NEXTLINE%
	-d "{\"country\":\"Russia\", \"year\": 1984}" http://localhost:8000/winners
	
	POSTing with invalid method:
	
	curl -i -X PUT -d "{\"country\":\"Russia\", \"year\": 2030}" http://localhost:8000/winners`

	usage = strings.ReplaceAll(usage, "%TOKEN%", AccessToken)
	if runtime.GOOS == "windows" {
		fmt.Println(strings.ReplaceAll(usage, "%NEXTLINE%", "^"))
	} else {
		fmt.Println(strings.ReplaceAll(usage, "%NEXTLINE%", "\\"))
	}
}
