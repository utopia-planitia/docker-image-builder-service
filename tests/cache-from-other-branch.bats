
load test_helper

setup() {
  export DATE=$(date +%s%N)

  export DOCKER_HOST=tcp://builder-0.builder:2375
  docker pull alpine:3.7 >&2
  docker build --build-arg version="$DATE" --cache-from currentBranch=$DATE -t cache-from-other-branch:$DATE tests/example-build >&2
  export DOCKER_HOST=tcp://builder-1.builder:2375
  docker pull alpine:3.7 >&2
}

@test "use cache from other branch" {
  run docker build --memory-swap=-1 --build-arg version="$DATE" --cache-from branch=$DATE -t cache-from-other-branch:$DATE tests/example-build
  [ "$status" -eq 0 ]
  [ "${lines[7]}" = " ---> Using cache" ]
  [ "${lines[10]}" = " ---> Using cache" ]
  [[ "${lines[12]}" =~ Successfully.* ]]
  [[ "${lines[13]}" =~ Successfully.* ]]
}
