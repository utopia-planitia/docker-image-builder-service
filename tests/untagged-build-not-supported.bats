
load test_helper

@test "local cache without tag" {
  export DATE=$(date +%s%N)
  docker pull alpine:3.7 >&2
  run docker build --memory-swap=-1 --build-arg version="$DATE" tests/example-build
  [ "$status" -ne 0 ]
  [ "${lines[1]}" = "Error response from daemon: failed to parse request: tag parameter is not set exactly once" ]
}
