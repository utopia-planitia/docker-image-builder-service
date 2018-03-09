
teardown () {
  echo teardown log
  echo "exit code: $status"
  for i in "${!lines[@]}"; do 
    printf "%s\t%s\n" "$i" "${lines[$i]}"
  done
  echo teardown done
}
