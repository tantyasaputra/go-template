#!/bin/bash
echo "Cleaning mocks folder"
find mocks/ -mindepth 1 -type d -exec rm -r {} + 2>/dev/null

echo "Run mockery"
mockery

echo "Move internal_ to mocks"
mv ./mocks/internal_/* ./mocks/ && rm -r ./mocks/internal_

echo "Finished"
