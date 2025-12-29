
if ! podman volume exists pgdata; then
	podman volume create pgdata
fi
podman create \
  --name postgres \
  -e POSTGRES_USER=test \
  -e POSTGRES_PASSWORD=test \
  -e POSTGRES_DB=testdb \
  -p 5432:5432 \
  -v pgdata:/var/lib/postgresql:Z \
  docker.io/library/postgres:latest

