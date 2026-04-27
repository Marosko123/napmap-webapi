param (
    $command
)

if (-not $command)  {
    $command = "start"
}

$ProjectRoot = "${PSScriptRoot}/.."

$env:NAPMAP_API_ENVIRONMENT="Development"
$env:NAPMAP_API_PORT="8080"
$env:NAPMAP_API_MONGODB_USERNAME="root"
$env:NAPMAP_API_MONGODB_PASSWORD="root"

function mongo {
    docker compose --file ${ProjectRoot}/deployments/docker-compose/compose.yaml $args
}

switch ($command) {
    "openapi" {
        docker run --rm -ti  -v ${ProjectRoot}:/local openapitools/openapi-generator-cli generate -c /local/scripts/generator-cfg.yaml
    }
    "start" {
        try {
            mongo up --detach
            go run ${ProjectRoot}/cmd/napmap-api-service
        } finally {
            mongo down
        }
    }
    "mongo" {
        mongo up
    }
    "docker" {
        docker build -t bednarmaros341/napmap-webapi:local-build -f ${ProjectRoot}/build/docker/Dockerfile .
    }
    "test" {
        go test -v ./...
    }
    default {
        throw "Unknown command: $command"
    }
}
