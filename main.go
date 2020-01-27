//package main
//
//import (
//	"golang.org/x/net/context"
//	"golang.org/x/tools/go/ssa/interp/testdata/src/os"
//	"log"
//	"os/exec"
//)

package main

import (
"context"
"github.com/google/go-github/github"
"go.uber.org/zap"
"golang.org/x/oauth2"
"log"
"os"
"os/exec"
)



var logger *zap.Logger

func InitLogger(e string) {
	switch e {
	case "production":
		logger, _ = zap.NewProduction()
		logger.Info("Using production level logging")
	case "development":
		logger, _ = zap.NewDevelopment()
		logger.Info("Using development level logging")
	default:
		logger = zap.NewExample()
		logger.Info("Using example level debugging")
	}
	defer logger.Sync()
	logger.Debug("Exiting func InitLogger")
}

func CloneRepos(github_personal_token string, github_org_name string) {
	logger.Debug("Entering func CloneRepos")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: github_personal_token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{Type: "Private"}

	// get all pages of results
	var allRepos []*github.Repository
	//PeopleNet
	for {
		logger.Info("Cloning every repo in Github Org", zap.String("github_org_name", github_org_name))
		repos, resp, err := client.Repositories.ListByOrg(ctx, github_org_name, opt)
		if err != nil {
			log.Fatal(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}

		//Found via https://yourbasic.org/golang/for-loop-range-array-slice-map-channel/
		for i := range repos {
			logger.Info("Cloning repo", zap.String("repo", repos[i].GetSSHURL()))
			cmd := exec.Command("git", "clone", repos[i].GetSSHURL())
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
		opt.Page = resp.NextPage
	}
	logger.Debug("Exiting func CloneRepos")
}

func main() {


	environment := os.Getenv("ENVIRONMENT")
	InitLogger(environment)

	logger.Debug("Entering func main")

	var githubPersonalToken = os.Getenv("GITHUB_PERSONAL_TOKEN")
	var githubOrgName = os.Getenv("GITHUB_ORG_NAME")

	logger.Info("GITHUB_ORG_NAME being used", zap.String("github_org_name", githubOrgName))
	logger.Debug("Logging initialized")

	CloneRepos(githubPersonalToken,githubOrgName)
	logger.Debug("Exiting func main")
}