
load test_helper

@test "uncached bytes should be 0 before the first build and should be 0 after a build" {

  run curl -H "GitBranchName: branchname-$DATE" --silent "http://builder-0.builder.container-image-builder.svc.cluster.local:2375/uncachedSize?cachefrom=%5B%5D&t=uncached-$DATE"
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "0" ]
  run curl -H "GitBranchName: branchname-$DATE" --silent "http://builder-1.builder.container-image-builder.svc.cluster.local:2375/uncachedSize?cachefrom=%5B%5D&t=uncached-$DATE"
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "0" ]

  # build on builder 1
  export DOCKER_HOST=tcp://builder-0.builder:2375
  docker pull alpine:3.7
  run docker build --build-arg version="$DATE" -t uncached-$DATE tests/example-build
  [ "$status" -eq 0 ]

  run curl -H "GitBranchName: branchname-$DATE" --silent "http://builder-0.builder.container-image-builder.svc.cluster.local:2375/uncachedSize?cachefrom=%5B%5D&t=uncached-$DATE"
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "0" ]
  run curl -H "GitBranchName: branchname-$DATE" --silent "http://builder-1.builder.container-image-builder.svc.cluster.local:2375/uncachedSize?cachefrom=%5B%5D&t=uncached-$DATE"
  [ "$status" -eq 0 ]
  [ "${lines[0]}" != "0" ]


  # build on builder 0
  export DOCKER_HOST=tcp://builder-1.builder:2375
  docker pull alpine:3.7
  run docker build --build-arg version="$DATE" -t uncached-$DATE tests/example-build
  [ "$status" -eq 0 ]
  [ "${lines[7]}" = " ---> Using cache" ]
  [ "${lines[10]}" = " ---> Using cache" ]
  [[ "${lines[12]}" =~ Successfully.* ]]
  [[ "${lines[13]}" =~ Successfully.* ]]

  run curl -H "GitBranchName: branchname-$DATE" --silent "http://builder-0.builder.container-image-builder.svc.cluster.local:2375/uncachedSize?cachefrom=%5B%5D&t=uncached-$DATE"
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "0" ]
  run curl -H "GitBranchName: branchname-$DATE" --silent "http://builder-1.builder.container-image-builder.svc.cluster.local:2375/uncachedSize?cachefrom=%5B%5D&t=uncached-$DATE"
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "0" ]
}
