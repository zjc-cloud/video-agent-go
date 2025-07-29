#!/bin/bash

# Global Git Commit & Push with Claude-generated messages
# Usage: gcp (in any git repository)

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Check if in git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

# Check for changes
if [[ -z $(git status -s) ]]; then
    echo -e "${YELLOW}No changes to commit${NC}"
    exit 0
fi

# Show status
echo -e "${GREEN}📋 Current changes:${NC}"
git status -s
echo ""

# Add all changes
git add -A

# Get diff for commit message generation
DIFF_CONTENT=$(git diff --cached --stat && echo "---" && git diff --cached | head -100)

# Generate commit message using Claude
echo -e "${GREEN}🤖 Generating commit message with Claude...${NC}"

# Create a temporary file for the response
TEMP_FILE=$(mktemp)

# Call Claude API
curl -s https://api.anthropic.com/v1/messages \
  -H "Content-Type: application/json" \
  -H "x-api-key: $ANTHROPIC_API_KEY" \
  -H "anthropic-version: 2023-06-01" \
  -d "{
    \"model\": \"claude-3-haiku-20240307\",
    \"max_tokens\": 100,
    \"messages\": [{
      \"role\": \"user\",
      \"content\": \"Based on this git diff, generate a concise commit message following conventional commits format (feat/fix/docs/style/refactor/test/chore). Only respond with the commit message, nothing else:\n\n$DIFF_CONTENT\"
    }]
  }" | jq -r '.content[0].text' > "$TEMP_FILE"

COMMIT_MSG=$(cat "$TEMP_FILE")
rm "$TEMP_FILE"

# Check if we got a valid response
if [ -z "$COMMIT_MSG" ] || [ "$COMMIT_MSG" == "null" ]; then
    echo -e "${RED}Failed to generate commit message. Using fallback.${NC}"
    COMMIT_MSG="chore: update files"
fi

# Show generated message
echo -e "${GREEN}📝 Commit message:${NC} $COMMIT_MSG"
echo ""

# Commit
git commit -m "$COMMIT_MSG"

if [ $? -ne 0 ]; then
    echo -e "${RED}Commit failed${NC}"
    exit 1
fi

# Show commit
echo -e "${GREEN}✅ Committed:${NC}"
git log -1 --oneline
echo ""

# Get current branch and check upstream
BRANCH=$(git branch --show-current)
HAS_UPSTREAM=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null)

# Determine push command
if [ -z "$HAS_UPSTREAM" ]; then
    PUSH_CMD="git push -u origin $BRANCH"
    echo -e "${YELLOW}⚠️  No upstream set. Will use: $PUSH_CMD${NC}"
else
    PUSH_CMD="git push"
fi

# Ask for push confirmation
echo -e "${YELLOW}Push to remote?${NC} [y/N] "
read -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${GREEN}🚀 Pushing...${NC}"
    eval $PUSH_CMD
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Push successful!${NC}"
    else
        echo -e "${RED}❌ Push failed${NC}"
        echo -e "${YELLOW}💡 Try: git pull --rebase origin $BRANCH${NC}"
    fi
else
    echo -e "${YELLOW}⏸️  Push cancelled. Changes committed locally.${NC}"
    echo -e "${YELLOW}To push later: $PUSH_CMD${NC}"
fi