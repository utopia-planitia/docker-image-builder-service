
setup() {
  export DATE=$(date +%s%N)
  docker system prune -af
  docker load -i alpine37.tar >&2
  docker build --build-arg version="$DATE" tests/example-build >&2
}

@test "build local cached image" {
  run docker build --build-arg version="$DATE" tests/example-build
  [ "$status" -eq 0 ]
  [ "${lines[4]}" = " ---> Using cache" ]
  [ "${lines[7]}" = " ---> Using cache" ]
  [[ "${lines[9]}" =~ Successfully.* ]]
}
