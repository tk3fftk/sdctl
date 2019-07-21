#!/bin/sh -e

go get github.com/linyows/git-semv/cmd/git-semv
go mod tidy

echo "latest tag is '$(git-semv latest)'"

latest_commit_log="$(git log -1 --pretty=format:'%s')"

if echo "$latest_commit_log" | grep -E '^fix' > /dev/null ; then
  tag=$(git-semv patch)
fi
if echo "$latest_commit_log" | grep -E '^feat' > /dev/null ; then
  tag=$(git-semv minor)
fi

if [ -n "$tag" ]; then
  echo "next tag is '$tag'"
  git tag "$tag"
  git push origin --tags
  # goreleaser needs GITHUB_TOKEN
  export GITHUB_TOKEN=$SCM_ACCESS_TOKEN
  curl -sL https://git.io/goreleaser | bash -s -- --rm-dist
else
  echo "no release will be published at this time"
fi

