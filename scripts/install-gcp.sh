#!/bin/bash

# Install gcp command globally

echo "Installing gcp command globally..."

# Create the global command
sudo tee /usr/local/bin/gcp > /dev/null << 'EOF'
#!/bin/bash

# Git Commit & Push with auto-generated message

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

# Show changes
echo "📋 Changes:"
git status -s
echo ""

# Add all changes
git add -A

# Generate commit message based on changes
ADDED=$(git diff --cached --name-status | grep "^A" | wc -l | tr -d ' ')
MODIFIED=$(git diff --cached --name-status | grep "^M" | wc -l | tr -d ' ')
DELETED=$(git diff --cached --name-status | grep "^D" | wc -l | tr -d ' ')

# Get first changed file to determine scope
FIRST_FILE=$(git diff --cached --name-only | head -1)
if [[ $FIRST_FILE == *"/"* ]]; then
    SCOPE=$(echo "$FIRST_FILE" | cut -d'/' -f1)
else
    SCOPE="project"
fi

# Determine commit type and message
if [ "$ADDED" -gt 0 ] && [ "$MODIFIED" -eq 0 ] && [ "$DELETED" -eq 0 ]; then
    TYPE="feat"
    MSG="add new files"
elif [ "$MODIFIED" -gt 0 ] && [ "$ADDED" -eq 0 ] && [ "$DELETED" -eq 0 ]; then
    TYPE="fix"
    MSG="update files"
elif [ "$DELETED" -gt 0 ] && [ "$ADDED" -eq 0 ] && [ "$MODIFIED" -eq 0 ]; then
    TYPE="refactor"
    MSG="remove files"
elif [ "$ADDED" -eq 0 ] && [ "$MODIFIED" -eq 0 ] && [ "$DELETED" -eq 0 ]; then
    TYPE="chore"
    MSG="update project"
else
    TYPE="chore"
    PARTS=()
    [ "$ADDED" -gt 0 ] && PARTS+=("add $ADDED")
    [ "$MODIFIED" -gt 0 ] && PARTS+=("update $MODIFIED")
    [ "$DELETED" -gt 0 ] && PARTS+=("remove $DELETED")
    MSG=$(IFS=", "; echo "${PARTS[*]}")
fi

COMMIT_MSG="$TYPE($SCOPE): $MSG

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>"

# Commit
echo "💬 Committing: $TYPE($SCOPE): $MSG"
git commit -m "$COMMIT_MSG"

if [ $? -ne 0 ]; then
    echo "❌ Commit failed"
    exit 1
fi

# Show commit
echo ""
echo "✅ Committed:"
git log -1 --oneline
echo ""

# Get branch info
BRANCH=$(git branch --show-current)
HAS_UPSTREAM=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null)

# Prepare push command
if [ -z "$HAS_UPSTREAM" ]; then
    PUSH_CMD="git push -u origin $BRANCH"
    echo "⚠️  No upstream branch. Will use: $PUSH_CMD"
else
    PUSH_CMD="git push"
fi

# Ask for push confirmation
echo ""
echo -n "🚀 Push to remote? [y/N] "
read -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Pushing..."
    eval $PUSH_CMD
    
    if [ $? -eq 0 ]; then
        echo "✅ Push successful!"
    else
        echo "❌ Push failed"
        echo "💡 Try: git pull --rebase origin $BRANCH"
        exit 1
    fi
else
    echo "⏸️  Cancelled. Commit saved locally."
    echo "To push later: $PUSH_CMD"
fi
EOF

# Make it executable
sudo chmod +x /usr/local/bin/gcp

echo "✅ Installation complete!"
echo ""
echo "Usage: gcp"
echo "This command will:"
echo "  1. Add all changes"
echo "  2. Generate a commit message automatically"
echo "  3. Commit the changes"
echo "  4. Ask for confirmation before pushing"
echo "  5. Handle first-time push with -u flag"