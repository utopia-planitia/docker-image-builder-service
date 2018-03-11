
load test_helper

@test "docker ps" {
  run docker ps
  [ "$status" -ne 0 ]
}

@test "docker run" {
  run docker run hello-world
  [ "$status" -ne 0 ]
}
