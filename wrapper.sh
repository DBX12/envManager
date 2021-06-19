#!/usr/bin/env bash


function envManager() {
  # get the first parameter (= verb like load, unload or debug)
  local verb=$1
  # get a temporary file where the stderr will be stored. The file has mode 600 by default.
  local TMPFILE=$(mktemp -t "XXXXXXXXXXXXXX")
  # call the binary and redirect its stderr into TMPFILE. If the binary is not in your PATH, put an absolute path here.
  envManager-bin "$@" 2> "$TMPFILE"
  # collect the exit code of the envManager-bin call
  local exitCode=$?
  # read and destroy the TMPFILE
  local tmpValue=$(cat "$TMPFILE"; rm -f "$TMPFILE")

  if ([[ $verb == "load" ]] || [[ $verb == "unload" ]]) && [[ $exitCode -eq 0 ]]; then
    # the output of these verbs should be eval'ed if the exitCode was zero
    eval "$tmpValue"
  elif [[ $verb != "__complete" ]]; then
    # for all other cases, just show what the binary returned in stderr
    echo "$tmpValue"
  fi
}