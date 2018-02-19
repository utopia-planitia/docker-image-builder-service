
setup() {
  docker load -i alpine37.tar >&2
  export DATE=$(date)
}

@test "build image" {
  run docker build --no-cache --build-arg version="$DATE" tests/example-build
  [ "$status" -eq 0 ]
  [[ "${lines[11]}" =~ Successfully.* ]]
}
