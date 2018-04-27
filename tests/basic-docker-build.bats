
load test_helper

@test "basic docker build" {
  docker pull alpine:3.7 >&2
  run docker build --memory-swap=-1 --build-arg version="$DATE" -t simple-build-$DATE tests/example-build
  [ "$status" -eq 0 ]
  [[ "${lines[12]}" =~ Successfully.* ]]
}
