
load test_helper

export DATE=$(date +%s%N)

@test "local cache without tag" {
  docker pull alpine:3.7 >&2
  docker build --build-arg version="$DATE" tests/example-build >&2

  run docker build --memory-swap=-1 --build-arg version="$DATE" tests/example-build
  [ "$status" -eq 0 ]
  [ "${lines[7]}" = " ---> Using cache" ]
  [ "${lines[10]}" = " ---> Using cache" ]
  [[ "${lines[12]}" =~ Successfully.* ]]
}

@test "local cache disabled with tag" {
  run docker build --memory-swap=-1 --build-arg version="$DATE" -t disaled-local-$DATE tests/example-build
  [ "$status" -eq 0 ]
  [ "${lines[7]}" = " ---> Using cache" ]
  [ "${lines[10]}" != " ---> Using cache" ]
  [[ "${lines[12]}" =~ Successfully.* ]]
  [[ "${lines[13]}" =~ Successfully.* ]]
}
