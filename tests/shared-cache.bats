
setup() {
  export DATE=$(date +%s%N)

  export DOCKER_HOST=tcp://builder_1:2375
  docker system prune -af
  docker load -i alpine37.tar >&2

  export DOCKER_HOST=tcp://builder_2:2375
  docker system prune -af
  docker load -i alpine37.tar >&2

  export DOCKER_HOST=tcp://builder_1:2375
  docker build --build-arg version="$DATE" -t test:latest tests/example-build >&2
  export DOCKER_HOST=tcp://builder_2:2375
}

@test "use shared build cache" {
  run docker build --build-arg version="$DATE" -t test:latest tests/example-build
  [ "$status" -eq 0 ]
  [ "${lines[4]}" = " ---> Using cache" ]
  [ "${lines[7]}" = " ---> Using cache" ]
  [[ "${lines[9]}" =~ Successfully.* ]]
  [[ "${lines[10]}" =~ Successfully.* ]]
}
