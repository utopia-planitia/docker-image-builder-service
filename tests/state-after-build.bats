
load test_helper

@test "registry state directly after image push" {
  export DATE=$(date +%s%N)
  docker build --memory-swap=-1 --build-arg version="$DATE" -t push-state-test-$DATE tests/example-build >&2
  run curl --fail -H Accept:application/vnd.docker.distribution.manifest.v2+json http://cache:5000/v2/push-state-test-$DATE/manifests/latest
  [ "$status" -eq 0 ]
}
