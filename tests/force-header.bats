
load test_helper

@test "force GitBranchName header" {
  [ ! -d $HOME/.docker ] || rm -r $HOME/.docker
  export DATE=$(date +%s%N)
  run docker build --memory-swap=-1 -t force-header tests/example-build
  [ "$status" -ne 0 ]
  [[ "${lines[1]}" = "Error response from daemon: failed to parse request: Branch not set via http header 'GitBranchName'" ]]
}
