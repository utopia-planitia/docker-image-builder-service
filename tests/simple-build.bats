
setup() {
  export DATE=$(date +%s%N)
  docker system prune -af
  docker load -i alpine37.tar >&2
}

@test "build image" {
  run docker build --build-arg version="$DATE" tests/example-build
  [ "$status" -eq 0 ]
  [[ "${lines[13]}" =~ Successfully.* ]]
}
