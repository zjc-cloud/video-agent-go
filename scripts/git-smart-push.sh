#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to generate commit message based on changes
generate_commit_message() {
    # Get the list of changed files
    changed_files=$(git diff --cached --name-only)
    
    if [ -z "$changed_files" ]; then
        changed_files=$(git diff --name-only)
    fi
    
    # Count different types of changes
    added=$(git diff --cached --name-status | grep "^A" | wc -l | tr -d ' ')
    modified=$(git diff --cached --name-status | grep "^M" | wc -l | tr -d ' ')
    deleted=$(git diff --cached --name-status | grep "^D" | wc -l | tr -d ' ')
    
    # Analyze the changes to determine commit type
    if [ "$added" -gt 0 ] && [ "$modified" -eq 0 ] && [ "$deleted" -eq 0 ]; then
        commit_type="feat"
        action="add"
    elif [ "$modified" -gt 0 ] && [ "$added" -eq 0 ] && [ "$deleted" -eq 0 ]; then
        commit_type="fix"
        action="update"
    elif [ "$deleted" -gt 0 ] && [ "$added" -eq 0 ] && [ "$modified" -eq 0 ]; then
        commit_type="refactor"
        action="remove"
    else
        commit_type="chore"
        action="update"
    fi
    
    # Get the most significant file or directory changed
    main_change=$(echo "$changed_files" | head -1)
    
    # Extract component/module from path
    if [[ $main_change == *"/"* ]]; then
        component=$(echo "$main_change" | cut -d'/' -f1)
    else
        component="project"
    fi
    
    # Generate descriptive message
    if [ "$added" -gt 0 ]; then
        desc_parts+=("$added new file(s)")
    fi
    if [ "$modified" -gt 0 ]; then
        desc_parts+=("$modified modification(s)")
    fi
    if [ "$deleted" -gt 0 ]; then
        desc_parts+=("$deleted deletion(s)")
    fi
    
    description=$(IFS=", "; echo "${desc_parts[*]}")
    
    # Create commit message
    echo "$commit_type($component): $action $description"
}

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

# Check for changes
if [[ -z $(git status -s) ]]; then
    echo -e "${YELLOW}No changes to commit${NC}"
    exit 0
fi

# Show current status
echo -e "${GREEN}=== Git Status ===${NC}"
git status -s
echo ""

# Add all changes
echo -e "${GREEN}Adding all changes...${NC}"
git add -A

# Generate commit message
commit_msg=$(generate_commit_message)
echo -e "${GREEN}Generated commit message:${NC} $commit_msg"
echo ""

# Ask for confirmation or custom message
echo -e "${YELLOW}Options:${NC}"
echo "1) Use generated message"
echo "2) Enter custom message"
echo "3) Cancel"
read -p "Choose option (1-3): " choice

case $choice in
    1)
        final_msg="$commit_msg"
        ;;
    2)
        read -p "Enter custom commit message: " final_msg
        ;;
    3)
        echo -e "${RED}Cancelled${NC}"
        git reset HEAD
        exit 0
        ;;
    *)
        echo -e "${RED}Invalid option${NC}"
        git reset HEAD
        exit 1
        ;;
esac

# Commit with the message
echo -e "${GREEN}Committing...${NC}"
git commit -m "$final_msg"

if [ $? -ne 0 ]; then
    echo -e "${RED}Commit failed${NC}"
    exit 1
fi

# Show the commit
echo -e "${GREEN}=== Commit Details ===${NC}"
git log -1 --oneline
echo ""

# Check current branch
current_branch=$(git branch --show-current)
echo -e "${GREEN}Current branch:${NC} $current_branch"

# Check if branch has upstream
has_upstream=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null)

# Prepare push command
if [ -z "$has_upstream" ]; then
    echo -e "${YELLOW}No upstream branch set. Will use: git push -u origin $current_branch${NC}"
    push_command="git push -u origin $current_branch"
else
    echo -e "${GREEN}Will push to existing upstream${NC}"
    push_command="git push"
fi

# Ask for push confirmation
echo ""
echo -e "${YELLOW}Ready to push. Do you want to proceed?${NC}"
echo "Command: $push_command"
read -p "Push now? (y/n): " confirm

if [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]]; then
    echo -e "${GREEN}Pushing...${NC}"
    eval $push_command
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Push successful!${NC}"
        
        # Show remote URL
        remote_url=$(git remote get-url origin 2>/dev/null)
        if [ ! -z "$remote_url" ]; then
            echo -e "${GREEN}Remote:${NC} $remote_url"
        fi
    else
        echo -e "${RED}Push failed. You may need to pull first or resolve conflicts.${NC}"
        echo -e "${YELLOW}Try: git pull --rebase origin $current_branch${NC}"
    fi
else
    echo -e "${YELLOW}Push cancelled. Your commit is saved locally.${NC}"
    echo -e "${YELLOW}To push later, run: $push_command${NC}"
fi