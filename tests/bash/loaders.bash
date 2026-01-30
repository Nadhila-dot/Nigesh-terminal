#!/bin/bash

frames=("■░░░" "█■░░" "░█■░" "░░█■" "░░░█" "■░░█" "█■░█" "░█■░" "░░█■")
while true; do
  for frame in "${frames[@]}"; do
    echo -ne "\r$frame Nigesh is processing.."
    sleep 0.1
  done
done

