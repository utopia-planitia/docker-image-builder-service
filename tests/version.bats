
load test_helper

@test "get docker client and server version" {
  run docker version
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "Client: Docker Engine - Community" ]
  [ "${lines[1]}" = " Version:           19.03.11" ]
  [ "${lines[2]}" = " API version:       1.40" ]
  [ "${lines[8]}" = "Server: Docker Engine - Community" ]
  [ "${lines[10]}" = "  Version:          19.03.11" ]
  [ "${lines[11]}" = "  API version:      1.40 (minimum version 1.12)" ]
}
