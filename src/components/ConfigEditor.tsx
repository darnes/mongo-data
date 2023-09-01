import * as React from 'react';
import { ChangeEvent } from 'react'

import { Input, SecretInput, SecretTextArea, Field, FieldSet } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { DataSourceOptions, SecureJsonData } from '../types';


interface Props extends DataSourcePluginOptionsEditorProps<DataSourceOptions> { }

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;

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


  const onPasswordReset = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        password: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        password: '',
      },
    });
  };
  const onPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    // not sure why I don't need to update secureJsonFields ... anyway
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        password: event.target.value,
      },
    });
  };
  const onSSLClientCertReset = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        sslClientCert: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        sslClientCert: '',
      },
    });
  };
  const onSSLClientCertChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    // not sure why I don't need to update secureJsonFields ... anyway
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        sslClientCert: event.target.value,
      },
    });
  };

  const { jsonData, secureJsonFields } = options;
  const secureJsonData = (options.secureJsonData || {}) as SecureJsonData;

  return (
    <FieldSet>
      {/* <InlineFieldRow> */}
      <Field label="Connection String" >
        <Input
          value={jsonData.mongoConnectionString || ''}
          placeholder="mongodb://username@host:port/?directConnection=true"
          onChange={onMongoConnectionStringChange}
          width={90}
        />
      </Field>
      <Field label="Password" >
        <SecretInput
          isConfigured={(secureJsonFields && secureJsonFields.password) as boolean}
          value={secureJsonData.password || ''}
          placeholder="password"
          width={90}
          onReset={onPasswordReset}
          onChange={onPasswordChange}
        />
      </Field>
      <Field label="SSL Client Certificate" >
        <SecretTextArea
          isConfigured={(secureJsonFields && secureJsonFields.sslClientCert) as boolean}
          cols={90} // of 200
          value={secureJsonData.sslClientCert || ''}
          placeholder="SSL Client  Certificate contents"
          onReset={onSSLClientCertReset}
          onChange={onSSLClientCertChange}
        />
      </Field>
      {/* </InlineFieldRow> */}
    </FieldSet>
  );
}
