
setup() {
  export DATE=$(date +%s%N)
  docker pull alpine:3.7 >&2
}

teardown () {
  if [ "$status" -eq 0 ]; then
    echo teardown log
    echo "status: $status"
    printf '%s\n' "${my_array[@]}"
    echo teardown done
  fi
}

@test "multiple tags" {
  run docker build --memory-swap=-1 --build-arg version="$DATE" -t a -t b tests/example-build
  [ "$status" -eq 0 ]
  [[ "${lines[13]}" =~ Successfully.* ]]
  [[ "${lines[14]}" =~ Successfully.* ]]
  [[ "${lines[15]}" =~ Successfully.* ]]
}
