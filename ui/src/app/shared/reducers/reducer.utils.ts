import { SerializedError, UnknownAction } from '@reduxjs/toolkit';
import { AxiosError, isAxiosError } from 'axios';

/**
 * Model for redux actions with pagination
 */
export type IQueryParams = { query?: string; page?: number; size?: number; sort?: string };

/**
 * Check if the async action type is rejected
 */
export function isRejectedAction(action: UnknownAction) {
  return action.type.endsWith('/rejected');
}

/**
 * Check if the async action type is pending
 */
export function isPendingAction(action: UnknownAction) {
  return action.type.endsWith('/pending');
}

/**
 * Check if the async action type is completed
 */
export function isFulfilledAction(action: UnknownAction) {
  return action.type.endsWith('/fulfilled');
}

const commonErrorProperties: Array<keyof SerializedError> = ['name', 'message', 'stack', 'code'];

/**
 * serialize function used for async action errors,
 * since the default function from Redux Toolkit strips useful info from axios errors
 */
export const serializeAxiosError = (value: unknown): AxiosError | SerializedError => {
  if (typeof value === 'object' && value !== null) {
    if (isAxiosError(value)) {
      return value;
    }
    const simpleError: SerializedError = {};
    for (const property of commonErrorProperties) {
      if (typeof value[property] === 'string') {
        simpleError[property] = value[property];
      }
    }

    return simpleError;
  }
  return { message: String(value) };
};

export interface EntityState<T> {
  loading: boolean;
  errorMessage: string | null;
  entities: ReadonlyArray<T>;
  entity: T;
  links?: unknown;
  updating: boolean;
  totalItems?: number;
  updateSuccess: boolean;
}
