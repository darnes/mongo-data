# usage

Run following commands:
```sh
# download plugin
wget https://github.com/darnes/mongo-data/releases/download/v1.0.1/ueon-mongodata-datasource-1.0.0.zip
#  extract  
tar xvf ueon-mongodata-datasource-1.0.0.zip

touch  docker-compose.yaml
```

Put following contents into docker-compose.yaml file, or just  [download it](./docker-compose.yaml)
```yaml
version: '3.1'

services:

  mongo:
    image: mongo
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: admin

  grafana:
    image: grafana/grafana-oss
    ports:
      - 3000:3000/tcp
    environment:
      GF_DATE_FORMATS_DEFAULT_TIMEZONE: "UTC+00:00"
      GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS: "ueon-mongodata-datasource"

    volumes:
      - ./ueon-mongodata-datasource:/var/lib/grafana/plugins/ueon-mongodata-datasource

```

## start containers  
```sh
docker compose up
```

## insert sample data to mongo:
get mongo DB container  ID:
```sh
docker ps | grep mongo
```

Get to mongo shell:
Use container ID from previous step instead of `bc6e35ca8698`
```sh
docker exec  -it bc6e35ca8698 -u admin -p admin
```

Run following commands to  create database, timeseries collection and insert some data

```js
use ts;
db.createCollection("stocks",{timeseries: {timeField: "date", metaField: "ticker", granularity: "seconds" }});
db.stocks.insertMany([
   { ticker: "MDB", date: ISODate("2021-12-18T15:59:00.000Z"), close: 252.47, volume: 55046.00},
   { ticker: "MDB", date: ISODate("2021-12-18T15:58:00.000Z"), close: 252.93, volume: 44042.00},
   { ticker: "MDB", date: ISODate("2021-12-18T15:57:00.000Z"), close: 253.61, volume: 40182.00},
   { ticker: "MDB", date: ISODate("2021-12-18T15:56:00.000Z"), close: 253.63, volume: 27890.00},
   { ticker: "MDB", date: ISODate("2021-12-18T15:55:00.000Z"), close: 254.03, volume: 40270.00}
]);
```

## Connect
Navigate to grafana: `http://localhost:3000/`, login, navigate to [datasources](http://localhost:3000/connections/datasources/new);
Select `Mondo Data`, use connection string: `mongodb://admin:admin@mongo:27017/`, click `Save & test`;

Make sure that you can see banner `connection to DB established`

Now, navigate to New Panel creation, enter following details:
DB: `ts`
collection: `stocks`
`Query`: `[]`
Click `Save Dashboard`;
You should be able to see the data now:
![Sample Panel](../img/sample-panel.png)

## Make it precise

- Go to dashboard settings and choose UTC timezone(`Coordinated Universal Time`)
- Edit panel, update query to `[{"$project": {"ticker": 1, "date": 1, "close": 1}}]`
- On the panel edit page, make sure that `Decimals` field is set to `2`

Now you should be able to see the data exactly as inserted:
![sample panel precise data](../img/sample-panel-precise.png)

## Celebrate
- Star this repo. That  will let author know that it is helpful so he can contribute more, making it better. 