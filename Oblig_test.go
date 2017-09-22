package main

import "testing"

func TestCreateURL(t *testing.T) { //Tests that the output of the createURL() is correct.
	inputURL := "/projectinfo/v1/github.com/repos/JonShard/RadIdea/"
	controlURL := "http://api.github.com/repos/JonShard/RadIdea"

	outputURL := CreateURL(inputURL)

	if outputURL != controlURL {
		t.Error("\nURL gave the wroung output:\n\tShould had been " + controlURL + "\n\tIt is " + outputURL)
	}
}
