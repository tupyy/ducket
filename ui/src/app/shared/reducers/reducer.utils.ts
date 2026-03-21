import { SerializedError } from '@reduxjs/toolkit';
import { AxiosError, isAxiosError } from 'axios';

const commonErrorProperties: Array<keyof SerializedError> = ['name', 'message', 'stack', 'code'];

export const serializeAxiosError = (value: unknown): AxiosError | SerializedError => {
  if (typeof value === 'object' && value !== null) {
    if (isAxiosError(value)) {
      return value;
    }
    const simpleError: SerializedError = {};
    for (const property of commonErrorProperties) {
      if (typeof (value as any)[property] === 'string') {
        simpleError[property] = (value as any)[property];
      }
    }
    return simpleError;
  }
  return { message: String(value) };
};
