
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

@test "always pull base image" {
  run docker build --memory-swap=-1 --build-arg version="$DATE" tests/example-build
  [ "$status" -eq 0 ]
  [[ "${lines[1]}" = "Step 1/3 : FROM alpine:3.7" ]]
  [[ "${lines[2]}" = "3.7: Pulling from library/alpine" ]]
  [[ "${lines[4]}" = "Status: Image is up to date for alpine:3.7" ]]
}
