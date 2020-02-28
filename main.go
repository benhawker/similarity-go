// Task 1.
// In the data folder of this repo there is a CSV file called reactions.csv.
// It contains real data corresponding to how users on Otta have reacted to
// (saved or discarded) jobs on the platform.

// The reaction data consists of four columns:
// user_id - the integer ID of the user who liked or disliked the job
// job_id - the integer ID of the job the user interacted with
// direction - whether the user liked (true) or disliked (false) the job
// time - the timestamp corresponding to when they reacted to the job

// The similarity score between two users is the number of jobs which they both like.
// Find the two users with the highest similarity.

// Task 2.
// In the data folder there is an additional CSV file called jobs.csv.
// It contains unique integer IDs for over 12,000 jobs, along with integer IDs
// for the job's associated company.

// The similarity score between two companies is the number of users who
// like at least one job at both companies. Using both the reactions.csv
// and jobs.csv data, find the two companies with the highest similarity score.
// Answer: [Enter the two company IDs & their similarity score here]

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Reaction struct {
	userId    int
	jobId     int
	liked     bool
	timestamp string
}

type HighestScore struct {
	entityType          string
	entityOne           int
	entityTwo           int
	numberOfSharedLikes int
}

// Represents map[jobId][]int{userId, userId2, ...}
// For calculating similarity score between companies
var userLikesByJob map[int][]int

// Represents map[userId][]int{jobId, jobId2, ...}
// For calculating similarity score between users
var likedJobsByUser map[int][]int

// Represents map[companyId][]int{jobId, jobId2, ...}
var jobsByCompany map[int][]int

// Represents map[companyId][]int{userId, userId2, ...}
var userLikesByCompany map[int][]int

const (
	reactionsFilePath = "data/reactions.csv"
	jobsFilePath      = "data/jobs.csv"
)

func init() {
	readAndMapReactions()
	readAndMapJobs()
}

func main() {
	findHighestSimilarityScore("user", likedJobsByUser)
	findHighestSimilarityScore("company", userLikesByCompany)
}

func findHighestSimilarityScore(entityType string, mappedData map[int][]int) {
	highestScore := HighestScore{}

	for keyId, associatedIds := range mappedData {
		for innerKeyId, innerAssociatedIds := range mappedData {
			// Do not compare like for like
			if keyId == innerKeyId {
				continue
			}

			// Do not double compare i.e. we will compare 1 vs 2, so no need to do 2 vs 1.
			if innerKeyId > keyId {
				continue
			}

			numberOfSharedLikes := numberOfSameElements(associatedIds, innerAssociatedIds)

			if numberOfSharedLikes > highestScore.numberOfSharedLikes {
				highestScore = HighestScore{
					entityType:          entityType,
					entityOne:           keyId,
					entityTwo:           innerKeyId,
					numberOfSharedLikes: numberOfSharedLikes,
				}
			}
		}
	}

	fmt.Printf("The highest similarity score is between %s %d and %s %d. They have %d shared likes.\n",
		highestScore.entityType,
		highestScore.entityOne,
		highestScore.entityType,
		highestScore.entityTwo,
		highestScore.numberOfSharedLikes)
}

func readAndMapJobs() {
	csvFile, err := os.Open(jobsFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("error reading all lines: %v", err)
	}

	jobsByCompany = make(map[int][]int)

	for i, line := range lines {
		if i == 0 {
			continue
		}

		jobId, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatal(err)
		}

		companyId, err := strconv.Atoi(line[1])
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := jobsByCompany[companyId]; ok {
			jobsByCompany[companyId] = appendIfMissing(jobsByCompany[companyId], jobId)
		} else {
			jobsByCompany[companyId] = []int{jobId}
		}
	}

	userLikesByCompany = make(map[int][]int)

	for companyId, jobs := range jobsByCompany {
		for _, jobId := range jobs {
			if _, ok := userLikesByCompany[companyId]; ok {
				userLikesByCompany[companyId] = appendIfMissing(userLikesByCompany[companyId], userLikesByJob[jobId])
			} else {
				userLikesByCompany[companyId] = userLikesByJob[jobId]
			}
		}

		// userLikesByCompany[companyId] = unique(userLikesByCompany[companyId])
	}
}

func readAndMapReactions() {
	csvFile, err := os.Open(reactionsFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("error reading all lines: %v", err)
	}

	likedJobsByUser = make(map[int][]int)
	userLikesByJob = make(map[int][]int)

	for i, line := range lines {
		if i == 0 {
			continue
		}

		userId, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatal(err)
		}

		jobId, err := strconv.Atoi(line[1])
		if err != nil {
			log.Fatal(err)
		}

		liked, err := strconv.ParseBool(line[2])
		if err != nil {
			log.Fatal(err)
		}

		// The similarity score between two users or companies is the number of jobs which
		// they both like hence we can only consider liked==true rows
		if liked == true {
			if _, ok := likedJobsByUser[userId]; ok {
				likedJobsByUser[userId] = appendIfMissing(likedJobsByUser[userId], jobId)
			} else {
				likedJobsByUser[userId] = []int{jobId}
			}

			if _, ok := userLikesByJob[jobId]; ok {
				userLikesByJob[jobId] = appendIfMissing(userLikesByJob[jobId], userId)
			} else {
				userLikesByJob[jobId] = []int{userId}
			}
		}
	}
}

func appendIfMissing(existing []int, i interface{}) []int {
    switch v := i.(type) {
	case int:
		for _, ele := range existing {
	        if ele == v {
	            return existing
	        }
	    }
	    return append(existing, v)
	case []int:
		itemsToAppend := []int{}
		for _, ele := range v {
	        if !contains(existing, ele) {
	            itemsToAppend = append(itemsToAppend, ele)
	        }
	    }
	    if len(itemsToAppend) > 0 {
	    	return append(existing, itemsToAppend...)
	    } else {
	    	return existing
		}
	default:
		return []int{}
	}

}

// returns a unique subset of the int slice provided.
func unique(input []int) []int {
	u := make([]int, 0, len(input))
	m := make(map[int]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func numberOfSameElements(a, b []int) int {
	counter := 0

	for _, v := range a {
		if contains(b, v) {
			counter++
		}
	}

	return counter
}

func contains(slice []int, element int) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
}
