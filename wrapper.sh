#!/usr/bin/env bash

function envManager() {
  local verb TMPFILE exitCode tmpValue
  # get the first parameter (= verb like load, unload or debug)
  verb=$1
  # get a temporary file where the stderr will be stored. The file has mode 600 by default.
  TMPFILE=$(mktemp -t "XXXXXXXXXXXXXX")
  # call the binary and redirect its stderr into TMPFILE. If the binary is not in your PATH, put an absolute path here.
  envManager-bin "$@" 2> "$TMPFILE"
  # collect the exit code of the envManager-bin call
  exitCode=$?
  # read and destroy the TMPFILE
  tmpValue=$(cat "$TMPFILE"; rm -f "$TMPFILE")

  # write the error message to stderr if the binary returned a non-zero exit code
  if [[ $exitCode -ne 0 ]]; then
    echo "$tmpValue" 1>&2
    return $exitCode
  fi

  case "$verb" in
    # the output of these verbs should be eval'ed
    load|unload)
      eval "$tmpValue"
      ;;
    # the stderr output of the __complete verbs is to be discarded
    __complete*)
      ;;
    # for all other cases, just show what the binary returned in stderr
    *)
      echo "$tmpValue" 1>&2
      ;;
  esac
}
