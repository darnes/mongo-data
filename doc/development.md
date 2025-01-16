# Development

## Install dependencies:
```sh
go get -u github.com/grafana/grafana-plugin-sdk-go
go mod tidy
```

## Build backend
```sh
# to run with docker compose with  `npm run server`
mage build:linux
```

## Build frontend
```sh
npm install
npm run build
npm run server
```

Keep in mind that plugin is running in the scope of docker, ie if the mongo is on localhost then use `host.docker.internal` as the host name.
ie: `mongodb://admin:admin@host.docker.internal:27017/?directConnection=true`


Pack 
```sh
ueon-mongodata-datasource

mv dist/ ueon-mongodata-datasource
zip ueon-mongodata-datasource-1.0.0.zip ueon-mongodata-datasource -r
```