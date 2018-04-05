
load test_helper

@test "docker ps is allowed" {
  run docker ps
  [ "$status" -eq 0 ]
}

@test "docker run is disabled" {
  run docker run hello-world
  [ "$status" -ne 0 ]
}
