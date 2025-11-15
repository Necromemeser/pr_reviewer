package handlers

import (
	"fmt"
	"net/http"
)

func CreatePullRequestHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Team created")
}

func MergePRHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Team info")
}

func ReassignReviewerHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Reassigned")
}
