default_registry('localhost:12345', host_from_cluster='k3d-registry.localhost:12345')

docker_build('web-client-api', '.', dockerfile='deploy/api.Dockerfile')
k8s_yaml('deploy/api.yaml')
k8s_yaml('deploy/api-svc.yaml')

k8s_resource(new_name='api-svc', objects=['api-svc'], resource_deps=['api'])