#!/usr/bin/env bash

for repo in $(cat repos.txt)
    do
        name=$(basename $repo);
        if [[ ! -d "repo-$(basename $repo)" ]]; then
            git clone -q $repo "repo-$(basename $repo)";
        else
            git pull
        fi
    done
