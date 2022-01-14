package chain_registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type CosmosGithubRegistry struct{}

func (c CosmosGithubRegistry) ListChains() ([]string, error) {
	client := github.NewClient(http.DefaultClient)
	var chains []string

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*5)
	tree, res, err := client.Git.GetTree(
		ctx,
		"cosmos",
		"chain-registry",
		"master",
		false)
	if err != nil || res.StatusCode != 200 {
		return chains, err
	}

	for _, entry := range tree.Entries {
		if *entry.Type == "tree" && !strings.HasPrefix(*entry.Path, ".") {
			chains = append(chains, *entry.Path)
		}
	}
	return chains, nil
}

func (c CosmosGithubRegistry) GetChain(name string) (ChainInfo, error) {
	client := github.NewClient(http.DefaultClient)

	chainFileName := path.Join(name, "chain.json")
	fileContent, _, res, err := client.Repositories.GetContents(
		context.Background(),
		"cosmos",
		"chain-registry",
		chainFileName,
		&github.RepositoryContentGetOptions{})
	if err != nil || res.StatusCode != 200 {
		return ChainInfo{}, errors.Wrap(err, fmt.Sprintf("error fetching %s", chainFileName))
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return ChainInfo{}, err
	}

	var result ChainInfo
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return ChainInfo{}, err
	}

	return result, nil
}

func (c CosmosGithubRegistry) SourceLink() string {
	return "https://github.com/cosmos/chain-registry"
}
