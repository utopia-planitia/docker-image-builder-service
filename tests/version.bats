

teardown () {
  if [ "$status" -eq 0 ]; then
    echo teardown log
    echo "status: $status"
    printf '%s\n' "${my_array[@]}"
    echo teardown done
  fi
}

@test "get docker client and server version" {
  run docker version
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "Client:" ]
  [ "${lines[1]}" = " Version:      17.11.0-ce" ]
  [ "${lines[7]}" = "Server:" ]
  [ "${lines[8]}" = " Version:      17.11.0-ce" ]
}
