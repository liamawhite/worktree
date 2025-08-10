#!/bin/sh

# Anything here will be ran in the root of a newly created worktree
git pull {{.Base}} {{.Branch}}