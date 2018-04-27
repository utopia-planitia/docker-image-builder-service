
setup () {
  export DATE=$(date +%s%N)
  mkdir -p $HOME/.docker
  echo "{\"HttpHeaders\": {\"GitBranchName\": \"branchname-$DATE\"}}" > $HOME/.docker/config.json
}

teardown () {
  echo teardown log
  echo "exit code: $status"
  for i in "${!lines[@]}"; do 
    printf "line %s:\t%s\n" "$i" "${lines[$i]}"
  done
  [ ! -d $HOME/.docker ] || rm -r $HOME/.docker
  echo teardown done
}
