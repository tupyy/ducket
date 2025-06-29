-- Direct SQL script to insert 100 sample transactions with tag associations
-- Usage: psql -h localhost -U postgres -d postgres -f tools/insert_sample_data.sql

-- First, let's insert some sample tags that could be applied to transactions
INSERT INTO tags (value) VALUES 
    ('food'), ('transport'), ('entertainment'), ('utilities'), ('shopping'), 
    ('income'), ('healthcare'), ('education'), ('travel'), ('subscription'),
    ('groceries'), ('gas'), ('restaurant'), ('salary'), ('freelance')
ON CONFLICT (value) DO NOTHING;

-- Insert sample rules with patterns for automatic tagging
INSERT INTO rules (id, pattern, created_at) VALUES 
    ('food_01', '(?i)(starbucks|mcdonalds|chipotle|subway|pizza)', now()),
    ('food_02', '(?i)(kroger|walmart.*grocery|target.*food)', now()),
    ('transport', '(?i)(uber|lyft|shell|gas station)', now()),
    ('utilities', '(?i)(electric|water|gas company|at&t|verizon|comcast)', now()),
    ('shopping', '(?i)(amazon|target|best buy|home depot)', now()),
    ('entertain', '(?i)(netflix|spotify|steam)', now()),
    ('income', '(?i)(salary|bonus|freelance|consulting)', now()),
    ('healthcare', '(?i)(cvs|walgreens|pharmacy)', now())
ON CONFLICT (id) DO NOTHING;

-- Associate rules with tags
INSERT INTO rules_tags (rule_id, tag) VALUES 
    ('food_01', 'food'),
    ('food_01', 'restaurant'),
    ('food_02', 'food'),
    ('food_02', 'groceries'),
    ('transport', 'transport'),
    ('transport', 'gas'),
    ('utilities', 'utilities'),
    ('shopping', 'shopping'),
    ('entertain', 'entertainment'),
    ('entertain', 'subscription'),
    ('income', 'income'),
    ('income', 'salary'),
    ('healthcare', 'healthcare')
ON CONFLICT (rule_id, tag) DO NOTHING;

-- Insert 100 sample transactions with random data
WITH sample_data AS (
    SELECT 
        generate_series(1, 100) as id,
        -- Random date between start of year and now
        (DATE_TRUNC('year', CURRENT_DATE) + 
         (RANDOM() * (CURRENT_DATE - DATE_TRUNC('year', CURRENT_DATE)::DATE)) * INTERVAL '1 day' +
         (RANDOM() * 24) * INTERVAL '1 hour')::TIMESTAMP as random_date,
        -- 70% debit, 30% credit
        CASE WHEN RANDOM() < 0.7 THEN 'debit' ELSE 'credit' END as transaction_kind
),
transaction_details AS (
    SELECT 
        id,
        random_date,
        transaction_kind,
        -- Amount based on transaction type
        CASE 
            WHEN transaction_kind = 'debit' THEN ROUND((RANDOM() * 495 + 5)::NUMERIC, 2)
            ELSE ROUND((RANDOM() * 4950 + 50)::NUMERIC, 2)
        END as amount,
        -- Generate realistic transaction content
        CASE 
            WHEN transaction_kind = 'debit' THEN
                (ARRAY['Purchase', 'Payment', 'Withdrawal', 'Subscription', 'Fee'])[FLOOR(RANDOM() * 5 + 1)] ||
                ' - ' ||
                (ARRAY['Amazon', 'Walmart', 'Target', 'Starbucks', 'McDonalds', 'Shell Gas Station', 
                       'Kroger', 'Home Depot', 'Best Buy', 'Netflix', 'Spotify', 'Uber', 'Lyft',
                       'CVS Pharmacy', 'Chipotle', 'Subway', 'AT&T', 'Comcast'])[FLOOR(RANDOM() * 18 + 1)] ||
                ' - Card ending in ' || LPAD(FLOOR(RANDOM() * 10000)::TEXT, 4, '0')
            ELSE
                (ARRAY['Deposit', 'Transfer', 'Salary', 'Bonus', 'Refund', 'Interest', 'Dividend'])[FLOOR(RANDOM() * 7 + 1)] ||
                ' - ' ||
                (ARRAY['Bank Transfer', 'Salary Deposit', 'Freelance Payment', 'Investment Return', 
                       'Tax Refund', 'Cashback Reward', 'Interest Payment', 'Dividend Payment',
                       'Rental Income', 'Consulting Fee', 'Bonus Payment'])[FLOOR(RANDOM() * 11 + 1)] ||
                ' - Account deposit'
        END as content
    FROM sample_data
),
inserted_transactions AS (
    INSERT INTO transactions (hash, date, kind, content, amount)
    SELECT 
        MD5(transaction_kind || random_date::TEXT || amount::TEXT || content) as hash,
        random_date,
        transaction_kind,
        content,
        amount
    FROM transaction_details
    RETURNING id, content
)
-- Associate transactions with tags based on content patterns
INSERT INTO transactions_tags (transaction_id, tag_id, rule_id)
SELECT DISTINCT
    t.id,
    rt.tag,
    rt.rule_id
FROM inserted_transactions t
JOIN rules_tags rt ON (
    -- Food transactions
    (rt.rule_id = 'food_01' AND t.content ~* '(starbucks|mcdonalds|chipotle|subway|pizza)') OR
    (rt.rule_id = 'food_02' AND t.content ~* '(kroger|walmart|target)') OR
    -- Transport transactions  
    (rt.rule_id = 'transport' AND t.content ~* '(uber|lyft|shell|gas station)') OR
    -- Utilities transactions
    (rt.rule_id = 'utilities' AND t.content ~* '(electric|water|gas company|at&t|verizon|comcast)') OR
    -- Shopping transactions
    (rt.rule_id = 'shopping' AND t.content ~* '(amazon|target|best buy|home depot)') OR
    -- Entertainment transactions
    (rt.rule_id = 'entertain' AND t.content ~* '(netflix|spotify|steam)') OR
    -- Income transactions
    (rt.rule_id = 'income' AND t.content ~* '(salary|bonus|freelance|consulting)') OR
    -- Healthcare transactions
    (rt.rule_id = 'healthcare' AND t.content ~* '(cvs|walgreens|pharmacy)')
);

-- Display summary of inserted data
SELECT 
    'Sample transactions with tags inserted successfully!' as message,
    COUNT(*) as total_transactions,
    SUM(CASE WHEN kind = 'debit' THEN 1 ELSE 0 END) as debit_count,
    SUM(CASE WHEN kind = 'credit' THEN 1 ELSE 0 END) as credit_count,
    MIN(date) as earliest_date,
    MAX(date) as latest_date,
    SUM(CASE WHEN kind = 'debit' THEN amount ELSE 0 END) as total_debits,
    SUM(CASE WHEN kind = 'credit' THEN amount ELSE 0 END) as total_credits
FROM transactions 
WHERE date >= DATE_TRUNC('year', CURRENT_DATE);

-- Display tag association summary
SELECT 
    'Tag Association Summary' as summary_type,
    COUNT(DISTINCT tt.transaction_id) as tagged_transactions,
    COUNT(*) as total_tag_associations,
    COUNT(DISTINCT tt.tag_id) as unique_tags_used
FROM transactions_tags tt
JOIN transactions t ON tt.transaction_id = t.id
WHERE t.date >= DATE_TRUNC('year', CURRENT_DATE);

-- Display breakdown by tag
SELECT 
    tt.tag_id as tag,
    tt.rule_id as rule,
    COUNT(*) as transaction_count,
    ROUND(AVG(t.amount), 2) as avg_amount
FROM transactions_tags tt
JOIN transactions t ON tt.transaction_id = t.id
WHERE t.date >= DATE_TRUNC('year', CURRENT_DATE)
GROUP BY tt.tag_id, tt.rule_id
ORDER BY transaction_count DESC; 