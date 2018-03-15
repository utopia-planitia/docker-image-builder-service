
load test_helper

setup() {
  export DATE=$(date +%s%N)
  docker pull alpine:3.7 >&2
}

@test "multiple tags" {
  run docker build --memory-swap=-1 --build-arg version="$DATE" -t multiple-tags-a -t multiple-tags-b tests/example-build
  [ "$status" -eq 0 ]
  [[ "${lines[12]}" =~ Successfully.* ]]
  [[ "${lines[13]}" =~ Successfully.* ]]
  [[ "${lines[14]}" =~ Successfully.* ]]
}
