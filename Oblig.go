//https://stackoverflow.com/questions/24299818/golang-how-to-decode-json-into-structs
//https://stackoverflow.com/questions/17156371/how-to-get-json-response-in-golang
// Fikk hjelp av bjoorn.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

//ProjectInfo struct that will contain all the data before it gets marshal-ed.
type ProjectInfo struct {
	Name      string   `json:"name"`
	Owner     string   `json:"owner"`
	Committer string   `josn:"committer"`
	Commits   int      `json:"commits"`
	Languages []string `json:"language"`
}

//CreateURL splits the input URL into its seperate parts.
func CreateURL(URLIn string) string {
	fmt.Printf("\nCreateURL():\n\tInputURL: " + URLIn)
	var parts []string
	parts = strings.Split(URLIn, "/")

	var site, user, repo string //Will store each part of the URL that will be used to GET the jsons.
	if len(parts) != 8 {        //Error if wrong num of parameters.
		return ""
	}

	site = parts[3]
	user = parts[5]
	repo = parts[6]
	outputURL := "http://api." + site + "/repos/" + user + "/" + repo

	fmt.Printf("\n\tOutputURL: " + outputURL)

	return outputURL //Rebuild the URL to it can be used with the git api.
}

func handlerProjectInfo(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	project := new(ProjectInfo) //Need somewhere to store the cherry-picking from all the repo-jsons.

	repoURL := CreateURL(r.URL.Path) // Get the base URL for using the git api.
	if repoURL == "" {

		fmt.Fprintf(w, http.StatusText(403)) //If it fails to create URL it means the amount of parameters was wrong. Bad request.
		return
	}

	contributorsURL := repoURL + "/contributors" //Create the URL that will be used to get jsons about the contributers.
	languagesURL := repoURL + "/languages"       //Create the URL that will be used to get jsons about the languages used.

	repoResponse, err := http.Get(repoURL) //Get the http response for the repo information.
	if err != nil {

		fmt.Fprintf(w, http.StatusText(403))
		return
	}
	responseLanguages, err := http.Get(languagesURL) //Get the http response for the language information.
	if err != nil {

		fmt.Fprintf(w, http.StatusText(403))
		return
	}
	contributorsResponse, err := http.Get(contributorsURL) //Get the http response for the contributor information.
	if err != nil {

		fmt.Fprintf(w, http.StatusText(403))
		return
	}

	var respRepo map[string]interface{}                              //Create an interface to decode the jsons into.
	json.NewDecoder(repoResponse.Body).Decode(&respRepo)             //Decode the http response.
	var contributors []interface{}                                   //Create an interface to decode the jsons into.
	json.NewDecoder(contributorsResponse.Body).Decode(&contributors) //Decode the http response.
	var respLanguages map[string]interface{}                         //Create an interface to decode the jsons into.
	json.NewDecoder(responseLanguages.Body).Decode(&respLanguages)   //Decode the http response.

	var ownerInterface = respRepo["owner"].(map[string]interface{}) //Take the nested structure in repo-json, and put it in a seperate interface.
	var topContributer = contributors[0].(map[string]interface{})   //Filter out the top contributer from the list of contributers.

	//A map of keys with corresponding strings
	languagesList := make([]string, 0, len(respLanguages)) //A list that will contain just the keys.

	for key := range respLanguages { //Takes the key in the language interface, which is the name of the language
		languagesList = append(languagesList, key) //and puts it in a string[]
	}

	project.Owner = ownerInterface["login"].(string) //Puts the values into the projectInfo struct:
	project.Name = respRepo["name"].(string)
	project.Committer = topContributer["login"].(string)
	project.Commits = int(topContributer["contributions"].(float64))
	project.Languages = languagesList

	res, err := json.Marshal(project) //Encodes the struct to json so it can be sent.
	if err != nil {
		fmt.Printf("\n\nThe Marshal failed.")
	}
	fmt.Fprintln(w, string(res)) //Prints json to webpage.

}

func main() {
	/*
		file, err := os.Open(".env") // For read access.
		if err != nil {
			fmt.Printf("\n\nThe Marshal failed.")
		}
	*/
	http.HandleFunc("/projectinfo/v1/", handlerProjectInfo)
	http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), nil) // Keep serving all requests that is recieved.
}
