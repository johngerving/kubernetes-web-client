load('ext://secret', 'secret_create_generic', 'secret_from_dict')
load('ext://helm_resource', 'helm_resource', 'helm_repo')

default_registry('localhost:12345', host_from_cluster='k3d-registry.localhost:12345')

helm_repo('bitnami', 'https://charts.bitnami.com/bitnami')

# Create backend API secret from .env file
secret_create_generic('backend-secret', from_env_file='backend/api/.env')
# Create frontend secret from .env file
secret_create_generic('frontend-secret', from_env_file='frontend/.env')

docker_build('web-client/api', 'backend/api', dockerfile='backend/api/Dockerfile')
docker_build('web-client/frontend', 'frontend', dockerfile='frontend/Dockerfile')

helm_resource('postgresql', 'bitnami/postgresql', resource_deps=['bitnami'])
k8s_yaml('deploy/api.yaml')
k8s_yaml('deploy/frontend.yaml')
# k8s_yaml('deploy/api-svc.yaml')

k8s_resource('api', port_forwards=8090)
k8s_resource('frontend', port_forwards=3000)
# k8s_resource(new_name='api-svc', objects=['api-svc'], resource_deps=['api'])