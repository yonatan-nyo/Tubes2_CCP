#!/bin/sh

# Run backend, scraper, and frontend server concurrently
concurrently \
  "./scraper" \
  "./backend"
