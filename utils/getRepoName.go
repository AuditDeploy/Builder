package utils

import "strings"

func GetRepoName(repoUrl string) string {
	repoSplit := strings.Split(repoUrl, "/")
	repo := repoSplit[len(repoSplit)-1]
	repoName := strings.Split(repo, ".")[0]

	return repoName
}
