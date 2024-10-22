load('ext://secret', 'secret_create_generic', 'secret_from_dict')
load('ext://helm_resource', 'helm_resource', 'helm_repo')

default_registry('localhost:12345', host_from_cluster='k3d-registry.localhost:12345')

helm_repo('bitnami', 'https://charts.bitnami.com/bitnami')

# Create backend API secret from .env file
secret_create_generic('backend-secret', from_env_file='backend/.env')
# Create frontend secret from .env file
secret_create_generic('frontend-secret', from_env_file='frontend/.env')

docker_build('web-client/api', 'backend', dockerfile='backend/Dockerfile')
docker_build('web-client/frontend', 'frontend', dockerfile='frontend/Dockerfile')

k8s_yaml('deploy/api.yaml')
k8s_yaml('deploy/frontend.yaml')
k8s_yaml('deploy/api-svc.yaml')
k8s_yaml('deploy/ingress.yaml')

k8s_resource(new_name='backend-secret', objects=['backend-secret'])
k8s_resource(new_name='frontend-secret', objects=['frontend-secret'])

helm_resource('postgresql', 'bitnami/postgresql', resource_deps=['bitnami'], port_forwards="5432:5432")
k8s_resource('api', resource_deps=['backend-secret', 'postgresql'])
# k8s_resource(new_name='api-svc', objects=['api-svc'], resource_deps=['api'])
k8s_resource('frontend', resource_deps=['frontend-secret'], port_forwards=3000)

# k8s_resource(new_name='api-ingress', objects=['api-ingress'])