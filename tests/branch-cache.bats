
setup() {
  export DATE=$(date +%s%N)

  export DOCKER_HOST=tcp://builder-0.builder:2375
  docker load -i alpine37.tar >&2

  export DOCKER_HOST=tcp://builder-1.builder:2375
  docker load -i alpine37.tar >&2

  export DOCKER_HOST=tcp://builder-0.builder:2375
  docker build --build-arg version="$DATE" --cache-from $DATE -t test:$DATE tests/example-build >&2
  export DOCKER_HOST=tcp://builder-1.builder:2375
}

@test "use shared branch build cache" {
  run docker build --memory-swap=-1 --build-arg version="$DATE" --cache-from $DATE -t test:$DATE tests/example-build
  [ "$status" -eq 0 ]
  [ "${lines[4]}" = " ---> Using cache" ]
  [ "${lines[7]}" = " ---> Using cache" ]
  [[ "${lines[9]}" =~ Successfully.* ]]
  [[ "${lines[10]}" =~ Successfully.* ]]
}
