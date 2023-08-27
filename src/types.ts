import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  queryText?: string;
  collection: string;
  db: string;
}

export const DEFAULT_QUERY: Partial<MyQuery> = {
  db: "local",
  collection: "startup_log",
  queryText: "[]" //consider following query: [{$project:{"hostname": 1, version: "$buildinfo.version"}}]
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  defaultDBName: string;
  mongoConnectionString: string
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {

}
