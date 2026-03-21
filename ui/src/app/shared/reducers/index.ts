import transactions from './transaction.reducer';
import rules from './rule.reducer';
import fileImport from './import.reducer';
import dashboard from './dashboard.reducer';

const rootReducer = {
  transactions,
  rules,
  fileImport,
  dashboard,
};

export default rootReducer;
