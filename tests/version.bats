
load test_helper

@test "get docker client and server version" {
  run docker version
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "Client:" ]
  [ "${lines[1]}" = " Version:	17.12.1-ce" ]
  [ "${lines[7]}" = "Server:" ]
  [ "${lines[9]}" = "  Version:	17.12.1-ce" ]
}
