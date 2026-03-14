#!/usr/bin/env bash
## CHANGELOG GENERATOR SCRIPT
## Generates changelog entries from merged pull requests
## Usage: ./changelog.sh <start_pr> <end_pr>
##   Example: ./changelog.sh 2700 2800

set -euo pipefail

# Check if gh command is available
if ! command -v gh &> /dev/null; then
    echo "Error: gh (GitHub CLI) is not installed"
    echo "Install it from: https://cli.github.com/"
    exit 1
fi

# Validate arguments
if [[ $# -ne 2 ]]; then
    echo "Usage: $0 <start_pr_number> <end_pr_number>"
    echo "  Generates changelog entries for merged PRs in the given range"
    echo ""
    echo "Example:"
    echo "  $0 2700 2800"
    exit 1
fi

START=$1
END=$2
REPO="skycoin/skycoin"

echo "Fetching merged PRs from #${START} to #${END} in ${REPO}..."
echo ""

# Iterate through PR numbers in reverse order (newest first)
for PR_NUM in $(seq "$END" -1 "$START"); do
    # Fetch PR info using gh CLI
    PR_INFO=$(gh pr view "$PR_NUM" --repo "$REPO" --json title,state,mergedAt 2>/dev/null || echo "")
    
    # Skip if PR doesn't exist or wasn't found
    if [[ -z "$PR_INFO" ]]; then
        continue
    fi
    
    # Parse JSON response
    STATE=$(echo "$PR_INFO" | jq -r '.state')
    MERGED_AT=$(echo "$PR_INFO" | jq -r '.mergedAt')
    TITLE=$(echo "$PR_INFO" | jq -r '.title')
    
    # Only include merged PRs
    if [[ "$STATE" == "MERGED" ]] && [[ "$MERGED_AT" != "null" ]]; then
        echo "- ${TITLE} [#${PR_NUM}](https://github.com/${REPO}/pull/${PR_NUM})"
    fi
done

echo ""
echo "Done. Copy the output above to your CHANGELOG.md file."
