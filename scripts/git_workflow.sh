#!/bin/bash

cd ../

# Function to check if we're in a git repository
check_git_repo() {
    if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
        echo "Error: Not a git repository!"
        exit 1
    fi
}

# Function to select repository
select_repository() {
    echo "Available repositories:"
    repos=($(find . -name ".git" -type d -exec dirname {} \;))
    if [ ${#repos[@]} -eq 0 ]; then
        echo "No Git repositories found in the current directory!"
        exit 1
    fi

    PS3="Select repository (enter number): "
    select repo in "${repos[@]}"; do
        if [ -n "$repo" ]; then
            cd "$repo"
            echo "Selected repository: $repo"
            break
        else
            echo "Invalid selection. Please try again."
        fi
    done
}

# Function to select/create branch
manage_branch() {
    echo "\nCurrent branch: $(git branch --show-current)"
    echo "\nBranch options:"
    echo "1) Use current branch"
    echo "2) Switch to existing branch"
    echo "3) Create new branch"
    
    read -p "Select option (1-3): " branch_option

    case $branch_option in
        1)
            echo "Continuing with current branch..."
            ;;
        2)
            echo "\nAvailable branches:"
            git branch
            read -p "Enter branch name to switch to: " branch_name
            git checkout "$branch_name" || exit 1
            ;;
        3)
            read -p "Enter new branch name: " new_branch
            git checkout -b "$new_branch" || exit 1
            ;;
        *)
            echo "Invalid option!"
            exit 1
            ;;
    esac
}

# Function to select files to add
select_files() {
    echo "\nUntracked and modified files:"
    git status -s

    echo "\nOptions:"
    echo "1) Stage all changes"
    echo "2) Select specific files"
    read -p "Select option (1-2): " file_option

    case $file_option in
        1)
            git add .
            ;;
        2)
            while true; do
                echo "\nCurrent status:"
                git status -s
                read -p "Enter file to stage (or 'done' to finish): " file
                
                if [ "$file" = "done" ]; then
                    break
                elif [ -e "$file" ]; then
                    git add "$file"
                    echo "Added: $file"
                else
                    echo "File not found: $file"
                fi
            done
            ;;
        *)
            echo "Invalid option!"
            exit 1
            ;;
    esac
}

# Function to create commit
create_commit() {
    while true; do
        read -p "\nEnter commit message: " commit_msg
        if [ -n "$commit_msg" ]; then
            git commit -m "$commit_msg" && break
        else
            echo "Commit message cannot be empty!"
        fi
    done
}

# Function to push changes
push_changes() {
    current_branch=$(git branch --show-current)
    echo "\nPushing to remote..."
    
    if ! git remote | grep -q .; then
        echo "No remote repository configured!"
        read -p "Enter remote repository URL: " remote_url
        git remote add origin "$remote_url"
    fi

    if git push origin "$current_branch"; then
        echo "Successfully pushed to remote!"
    else
        echo "Failed to push to remote. Please check your credentials and try again."
        exit 1
    fi
}

# Main script execution
echo "=== Git Workflow Script ==="

# Select repository
select_repository

# Check if it's a git repository
check_git_repo

# Manage branch
manage_branch

# Select and add files
select_files

# Show status before commit
echo "\nCurrent status:"
git status

# Create commit
create_commit

# Push changes
read -p "\nDo you want to push changes? (y/n): " push_option
if [ "$push_option" = "y" ] || [ "$push_option" = "Y" ]; then
    push_changes
fi

echo "\nGit workflow completed successfully!"