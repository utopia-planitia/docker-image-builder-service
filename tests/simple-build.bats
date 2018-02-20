
setup() {
  export DATE=$(date +%s%N)
  docker load -i alpine37.tar >&2
}

@test "build image" {
  run docker build --memory-swap=-1 --build-arg version="$DATE" tests/example-build
  [ "$status" -eq 0 ]
  [[ "${lines[10]}" =~ Successfully.* ]]
}
