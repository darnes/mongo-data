import React, { ChangeEvent } from 'react';
import { InlineField, Input, TextArea } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { DataSourceOptions, MyQuery } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, DataSourceOptions>;

const isValidJson = function (text: string): boolean {
  try {
    JSON.parse(text);
    return true;
  } catch (exc) {
    return false;
  }
};

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const onQueryTextChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    // tbd: will be nice to auto-forma the query... may be later
    onChange({ ...query, queryText: event.target.value });

    // I'm guessing there might be better place for this check.
    if (isValidJson(event.target.value)) {
      onRunQuery();
    }
  };


  const onCollectionChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, collection: event.target.value });
    onRunQuery();
  };

  const onDBChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, db: event.target.value });
    onRunQuery();
  };

  const { queryText, collection, db } = query;

  const queryInvalid = !isValidJson(queryText || '{}');
  const queryError = queryInvalid ? "Query must be a valid JSON" : null;

  return (
    <div className="gf-form">
      <InlineField label="DB">
        <Input onChange={onDBChange} value={db} width={8} type="string" />
      </InlineField>
      <InlineField label="Collection">
        <Input onChange={onCollectionChange} value={collection} width={8} type="string" />
      </InlineField>
      {/*TBD:  consider making query input bigger */}
      <InlineField label="Query" invalid={queryInvalid} error={queryError}>
        <TextArea onChange={onQueryTextChange} value={queryText || ''} />
      </InlineField>
    </div>
  );
}
