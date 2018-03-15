
load test_helper

setup() {
  export DATE=$(date +%s%N)
  docker pull alpine:3.7 >&2
}

@test "basic docker build" {
  run docker build --memory-swap=-1 --build-arg version="$DATE" -t simple-build-$DATE tests/example-build
  [ "$status" -eq 0 ]
  [[ "${lines[12]}" =~ Successfully.* ]]
}
