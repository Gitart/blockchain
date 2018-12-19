#!/bin/bash

nc="\\033[0m"
red="\\033[31m"
green="\\033[32m"
yellow="\\033[33m"
blue="\\033[34m"
purple="\\033[35m"
cyan="\\033[36m"
white="\\033[37m"
bold="$(tput bold)"
normal="$(tput sgr0)"

printf "$red$bold%s$normal\\n" "Error: path is required"
printf "$green$bold%s$normal\\n" "Error: path is required"
printf "$yellow$bold%s$normal\\n" "Error: path is required"
printf "$blue%s$normal\\n" "Error: path is required"
printf "$purple$bold%s$normal\\n" "Error: path is required"
printf "$nc$bold%s$normal\\n" "Error: path is required"
printf "$cyan$bold%s$normal\\n" "Error: path is required"
printf "$white$bold%s$normal $red$bold%s$normal   \\n" "Error: path is required" "Repo34"
printf "$white$bold%s$normal $red$bold%s$normal  -- $green%s$normal \\n" "Error: path is required" "Repo34" "Test"

