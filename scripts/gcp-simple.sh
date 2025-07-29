#!/bin/bash

# Simple Git Commit & Push
# Usage: gcp-simple (in any git repository)

# Check if in git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "Error: Not in a git repository"
    exit 1
fi

# Check for changes
if [[ -z $(git status -s) ]]; then
    echo "No changes to commit"
    exit 0
fi

# Show status briefly
echo "Changes to commit:"
git status -s
echo ""

# Add all changes
git add -A

# Commit with generated message
git commit -m "🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>"

# Get current branch
BRANCH=$(git branch --show-current)

# Check if upstream exists
if ! git rev-parse --abbrev-ref --symbolic-full-name @{u} > /dev/null 2>&1; then
    PUSH_CMD="git push -u origin $BRANCH"
else  
    PUSH_CMD="git push"
fi

# Confirm push
echo "Push to remote? [y/N] "
read -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    $PUSH_CMD
else
    echo "Cancelled. To push later: $PUSH_CMD"
fi