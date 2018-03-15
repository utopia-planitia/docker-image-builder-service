
load test_helper

export DATE=$(date +%s%N)

@test "local cache without tag" {
  docker pull alpine:3.7 >&2
  run docker build --memory-swap=-1 --build-arg version="$DATE" tests/example-build
  [ "$status" -ne 0 ]
  [ "${lines[1]}" = "Error response from daemon: untagged builds are not supported" ]
}
