# How to use

## Connect 

Specify connection string as you would do for the API.
ie: `mongodb://username:password@you-host:port/?serverSelectionTimeoutMS=2000`
Make sure that you click save & test and don't see any errors.

## Query for data
Specify database and collection names in the appropriate fields. 
## Query
The query is the json version  of the query you can use in the aggregate method formatted as json. See [docs](https://www.mongodb.com/docs/manual/aggregation/) for details on aggregate method.

The query need to be a valid json, so you'll need to user modifiers ie `$match` quoted ie `"$match"`.

For example: 
to run query
```bson
[{$project:  {"_id":0}}]
```
must be written as:
```json
[{"$project":  {"_id":0}}]
```

Currently plugin allows you to render only 1st level data of the response. So to show nested values use projection. ie:
```json
[{"$project":  {"_id":0, "name": "$details.name"}}]
```
![Screen](./img/screen.png)  
