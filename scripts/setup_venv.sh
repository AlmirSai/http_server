#!/bin/bash

# Create virtual environment
python3 -m venv .venv

# Activate virtual environment
source .venv/bin/activate

# Install requirements
pip install -r requirements.txt

# Run linting on k8s_manager.py
echo "Running linters on k8s_manager.py..."
pylint scripts/k8s_manager.py
flake8 scripts/k8s_manager.py
black scripts/k8s_manager.py

echo "Setup complete! Virtual environment is ready."
echo "To activate the virtual environment, run: source .venv/bin/activate"