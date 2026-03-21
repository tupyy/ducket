import transactions from './transaction.reducer';
import rules from './rule.reducer';
import fileImport from './import.reducer';

const rootReducer = {
  transactions,
  rules,
  fileImport,
};

export default rootReducer;
