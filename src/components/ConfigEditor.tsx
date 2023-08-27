import React, { ChangeEvent } from 'react';
// import { InlineField, Input, SecretInput } from '@grafana/ui';
import { InlineField, Input } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
// import { MyDataSourceOptions, MySecureJsonData } from '../types';
import { MyDataSourceOptions } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> { }

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;

  // const onPathChange = (event: ChangeEvent<HTMLInputElement>) => {
  //   const jsonData = {
  //     ...options.jsonData,
  //     path: event.target.value,
  //   };
  //   onOptionsChange({ ...options, jsonData });
  // };

  // Secure field (only sent to the backend)
  const onMongoConnectionStringChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      mongoConnectionString: event.target.value,
    };
    onOptionsChange({
      ...options,
      jsonData,
    });
  };


  // const onResetMongoConnectionString = () => {
  //   onOptionsChange({
  //     ...options,
  //     secureJsonFields: {
  //       ...options.secureJsonFields,
  //       mongoConnectionString: false,
  //     },
  //     secureJsonData: {
  //       ...options.secureJsonData,
  //       mongoConnectionString: '',
  //     },
  //   });
  // };
  const { jsonData } = options;
  // const { jsonData, secureJsonFields } = options;
  // const { secureJsonFields } = options;
  // const secureJsonData = (options.secureJsonData || {}) as MySecureJsonData;

  return (
    <div className="gf-form-group">
      {/* <InlineField label="Path" labelWidth={12}>
        <Input
          onChange={onPathChange}
          value={jsonData.path || ''}
          placeholder="json field returned to frontend"
          width={40}
        />
      </InlineField>*/}
      {/* mongodb://mongo:27017/?directConnection=true&serverSelectionTimeoutMS=2000 */}
      <InlineField label="Connection String" labelWidth={12}>
        <Input
          value={jsonData.mongoConnectionString || ''}
          placeholder="mongodb://username:password@host:port/?directConnection=true&serverSelectionTimeoutMS=2000"
          width={100}
          onChange={onMongoConnectionStringChange}
        />
      </InlineField>
      {/* TBD find sort of help tool tip: Mongo DB connection string(Grafana backend Should be able to use it to connect) */}
      {/* <InlineField label="Connection String" labelWidth={20} >
        <SecretInput
          isConfigured={(secureJsonFields && secureJsonFields.mongoConnectionString) as boolean}
          value={secureJsonData.mongoConnectionString || ''}
          placeholder="mongodb://mongo:27017/?directConnection=true&serverSelectionTimeoutMS=2000"
          width={40}
          onReset={onResetMongoConnectionString}
          onChange={onMongoConnectionStringChange}
        />
      </InlineField> */}
    </div>
  );
}
