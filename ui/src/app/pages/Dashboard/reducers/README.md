# Dashboard Reducers

This folder contains dashboard-related reducers that **DO NOT** make backend requests. These reducers manage local state and calculations for the dashboard components.

## Reducers

### `label-report.reducer.ts`
- Calculates transaction amounts grouped by labels (key:value format)
- Processes transaction type reports (debit/credit totals)
- Generates account transaction type reports
- **Input**: Array of transactions (already fetched from backend)
- **Output**: Aggregated label amounts, transaction type summaries

### `monthly-label-report.reducer.ts`
- Calculates monthly trends for label amounts over time
- Groups transactions by month and label combinations
- Generates time-series data for charts and reports
- **Input**: Array of transactions (already fetched from backend)  
- **Output**: Monthly label totals, trend data for visualization

## Usage

Import actions and reducers from the index file:

```typescript
import { 
  calculateLabelReport, 
  calculateMonthlyLabelReport,
  labelReportReducer 
} from './reducers';
```

## Data Flow

1. Dashboard fetches transactions from backend (via shared `transaction.reducer.ts`)
2. Local reducers process the transaction data to generate reports
3. Charts and visualizations consume the processed data

## Note

These reducers use `createAsyncThunk` for async processing of large datasets, but they don't make HTTP requests. They only perform local calculations on the provided transaction data. 