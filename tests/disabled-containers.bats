
load test_helper

@test "docker ps is disabled" {
  run docker ps
  [ "$status" -ne 0 ]
}

@test "docker run is disabled" {
  run docker run hello-world
  [ "$status" -ne 0 ]
}
