#!/bin/sh

# Run backend, scraper, and frontend server concurrently
concurrently \
  "./backend" \
  "./scraper" \
  "serve -s ./frontend -l 4001"
