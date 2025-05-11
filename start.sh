#!/bin/sh

# Run backend, scraper, and frontend server concurrently
concurrently \
  "./scraper" \
  "./backend" \
  "serve -s ./frontend -l 4001"
