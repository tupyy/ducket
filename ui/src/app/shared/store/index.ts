import { AnyAction, configureStore, ThunkAction } from '@reduxjs/toolkit';
import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux';

import reducer from '@app/shared/reducers';
import loggerMiddleware from './logger-middleware';

const store = configureStore({
  reducer: reducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        // Ignore these field paths in all actions
        ignoredActionPaths: ['payload.config', 'payload.request', 'payload.headers', 'error', 'meta.arg', 'meta'],
        // Ignore these paths in the state
        ignoredPaths: ['transactions.transactions', 'transactionFilter.filteredTransactions', 'transactionFilter.sourceTransactions'],
        // Ignore specific action types that contain non-serializable data
        ignoredActions: ['transactionFilter/applyFilters/fulfilled', 'transactionFilter/applyFilters/pending', 'transactionFilter/applyFilters/rejected'],
        // Allow Date objects
        isSerializable: (value: any) => {
          return value instanceof Date || typeof value !== 'object' || value === null || Array.isArray(value) || typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean';
        },
      },
    }).concat(loggerMiddleware),
});

const getStore = () => store;

export type IRootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export const useAppSelector: TypedUseSelectorHook<IRootState> = useSelector;
export const useAppDispatch = () => useDispatch<AppDispatch>();
export type AppThunk<ReturnType = void> = ThunkAction<ReturnType, IRootState, unknown, AnyAction>;

export default getStore;
