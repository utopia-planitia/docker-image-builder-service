
disable_snapshots()
allow_k8s_contexts(['test', 'ci'])

k8s_yaml('kubernetes/builders.yaml')
k8s_yaml('kubernetes/cache.yaml')
k8s_yaml('kubernetes/devtools.yaml')
k8s_yaml('kubernetes/dispatcher.yaml')
k8s_yaml('kubernetes/mirror.yaml')
k8s_yaml('kubernetes/tests.yaml')

k8s_resource(
  'tests',
  trigger_mode=TRIGGER_MODE_MANUAL,
  resource_deps=['dispatcher', 'builder', 'cache', 'mirror'],
)

k8s_resource('dispatcher', port_forwards=['2375'])

docker_build(
  'dispatcher-image',
  './dispatcher',
  dockerfile='./dispatcher/Dockerfile',
)

docker_build(
  'worker-image',
  './worker',
  dockerfile='./worker/Dockerfile',
)

docker_build(
  'devtools-image',
  '.',
  dockerfile='./devtools/Dockerfile',
  only=['./devtools', './tests'],
)


