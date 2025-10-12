#!/bin/bash
# Test multi-file diff navigation

echo "Generating multi-file diff..."
(diff -u longtest1.txt longtest2.txt; diff -u test1.txt test2.txt; diff -u file1.1.go file1.2.go) | ./differential --navigate
