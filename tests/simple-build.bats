
setup() {
  export DATE=$(date +%s%N)
  docker pull alpine:3.7
}

@test "build image" {
  run docker build --memory-swap=-1 --build-arg version="$DATE" tests/example-build
  [ "$status" -eq 0 ]
  [[ "${lines[13]}" =~ Successfully.* ]]
}
