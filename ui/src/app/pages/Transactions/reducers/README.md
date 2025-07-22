# Transaction Reducers

This folder contains transaction-related reducers that **DO NOT** make backend requests. These reducers manage local state for the transaction components.

## Reducers

### `transaction-filter.reducer.ts`
- Manages filtering, sorting, pagination, and expansion state for transactions
- Handles date range selection
- Manages local UI state for the transaction list

### `transactionSummary.reducer.ts`
- Calculates and manages transaction summary data
- Aggregates totals by label keys or filtered transaction totals
- Provides summary statistics for the transaction view

## Usage

Import actions and reducers from the index file:

```typescript
import { 
  setSelectedLabels, 
  calculateTransactionSummary,
  transactionFilterReducer 
} from './reducers';
```

## Note

Reducers that make backend requests (like `transaction.reducer.ts`) remain in the shared reducers folder since they might be used across multiple components. 