-- Sample Transaction Data Generator
-- This procedure inserts 100 sample transactions with dates between the beginning of the year and now
-- Also creates sample rules and associates transactions with tags

CREATE OR REPLACE FUNCTION insert_sample_transactions()
RETURNS void AS $$
DECLARE
    i INTEGER;
    random_date TIMESTAMP;
    transaction_kind VARCHAR(30);
    transaction_content TEXT;
    transaction_amount NUMERIC(15,2);
    transaction_hash VARCHAR(100);
    transaction_id INTEGER;
    year_start DATE;
    current_date_val DATE;
    days_in_range INTEGER;
    
    -- Sample merchant names and transaction types
    merchants TEXT[] := ARRAY[
        'Amazon', 'Walmart', 'Target', 'Starbucks', 'McDonalds', 
        'Shell Gas Station', 'Kroger', 'Home Depot', 'Best Buy', 'Netflix',
        'Spotify', 'Uber', 'Lyft', 'Airbnb', 'PayPal', 'Apple Store',
        'Google Play', 'Steam', 'Adobe', 'Microsoft', 'Costco', 'CVS Pharmacy',
        'Walgreens', 'Chipotle', 'Subway', 'Pizza Hut', 'Dominos',
        'AT&T', 'Verizon', 'T-Mobile', 'Comcast', 'Electric Company',
        'Water Department', 'Gas Company', 'Insurance Co', 'Bank Transfer',
        'Salary Deposit', 'Freelance Payment', 'Investment Return', 'Tax Refund',
        'Gift Card', 'Cashback Reward', 'Interest Payment', 'Dividend Payment',
        'Rental Income', 'Side Hustle', 'Garage Sale', 'Online Sale',
        'Consulting Fee', 'Bonus Payment'
    ];
    
    transaction_types TEXT[] := ARRAY[
        'Purchase', 'Payment', 'Transfer', 'Deposit', 'Withdrawal', 
        'Subscription', 'Refund', 'Fee', 'Interest', 'Dividend',
        'Salary', 'Bonus', 'Commission', 'Rental', 'Utility'
    ];
    
BEGIN
    -- Insert sample tags
    INSERT INTO tags (value) VALUES 
        ('food'), ('transport'), ('entertainment'), ('utilities'), ('shopping'), 
        ('income'), ('healthcare'), ('education'), ('travel'), ('subscription'),
        ('groceries'), ('gas'), ('restaurant'), ('salary'), ('freelance')
    ON CONFLICT (value) DO NOTHING;
    
    -- Insert sample rules with patterns
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
    
    -- Calculate date range
    year_start := DATE_TRUNC('year', CURRENT_DATE);
    current_date_val := CURRENT_DATE;
    days_in_range := current_date_val - year_start;
    
    -- Insert 100 sample transactions
    FOR i IN 1..100 LOOP
        -- Generate random date between start of year and now
        random_date := year_start + (RANDOM() * days_in_range) * INTERVAL '1 day' 
                      + (RANDOM() * 24) * INTERVAL '1 hour'
                      + (RANDOM() * 60) * INTERVAL '1 minute'
                      + (RANDOM() * 60) * INTERVAL '1 second';
        
        -- Randomly choose debit (70%) or credit (30%) - more debits are realistic
        IF RANDOM() < 0.7 THEN
            transaction_kind := 'debit';
            -- Debit amounts typically smaller, between $5 and $500
            transaction_amount := ROUND((RANDOM() * 495 + 5)::NUMERIC, 2);
        ELSE
            transaction_kind := 'credit';
            -- Credit amounts can be larger (salary, etc), between $50 and $5000
            transaction_amount := ROUND((RANDOM() * 4950 + 50)::NUMERIC, 2);
        END IF;
        
        -- Generate realistic transaction content
        transaction_content := transaction_types[1 + FLOOR(RANDOM() * array_length(transaction_types, 1))] 
                             || ' - ' 
                             || merchants[1 + FLOOR(RANDOM() * array_length(merchants, 1))];
        
        -- Add some additional context for certain types
        IF transaction_kind = 'debit' THEN
            transaction_content := transaction_content || ' - Card ending in ' || LPAD(FLOOR(RANDOM() * 10000)::TEXT, 4, '0');
        ELSE
            transaction_content := transaction_content || ' - Account deposit';
        END IF;
        
        -- Generate hash similar to how the Go code does it (simplified version)
        transaction_hash := MD5(transaction_kind || random_date::TEXT || transaction_amount::TEXT || transaction_content);
        
        -- Insert the transaction and get the ID
        INSERT INTO transactions (hash, date, kind, content, amount)
        VALUES (transaction_hash, random_date, transaction_kind, transaction_content, transaction_amount)
        RETURNING id INTO transaction_id;
        
        -- Associate transaction with appropriate tags based on content
        -- Food transactions
        IF transaction_content ~* '(starbucks|mcdonalds|chipotle|subway|pizza)' THEN
            INSERT INTO transactions_tags (transaction_id, tag_id, rule_id) VALUES 
                (transaction_id, 'food', 'food_01'),
                (transaction_id, 'restaurant', 'food_01');
        ELSIF transaction_content ~* '(kroger|walmart|target)' THEN
            INSERT INTO transactions_tags (transaction_id, tag_id, rule_id) VALUES 
                (transaction_id, 'food', 'food_02'),
                (transaction_id, 'groceries', 'food_02');
        END IF;
        
        -- Transport transactions
        IF transaction_content ~* '(uber|lyft|shell|gas station)' THEN
            INSERT INTO transactions_tags (transaction_id, tag_id, rule_id) VALUES 
                (transaction_id, 'transport', 'transport'),
                (transaction_id, 'gas', 'transport');
        END IF;
        
        -- Utilities transactions
        IF transaction_content ~* '(electric|water|gas company|at&t|verizon|comcast)' THEN
            INSERT INTO transactions_tags (transaction_id, tag_id, rule_id) VALUES 
                (transaction_id, 'utilities', 'utilities');
        END IF;
        
        -- Shopping transactions
        IF transaction_content ~* '(amazon|target|best buy|home depot)' THEN
            INSERT INTO transactions_tags (transaction_id, tag_id, rule_id) VALUES 
                (transaction_id, 'shopping', 'shopping');
        END IF;
        
        -- Entertainment transactions
        IF transaction_content ~* '(netflix|spotify|steam)' THEN
            INSERT INTO transactions_tags (transaction_id, tag_id, rule_id) VALUES 
                (transaction_id, 'entertainment', 'entertain'),
                (transaction_id, 'subscription', 'entertain');
        END IF;
        
        -- Income transactions
        IF transaction_content ~* '(salary|bonus|freelance|consulting)' THEN
            INSERT INTO transactions_tags (transaction_id, tag_id, rule_id) VALUES 
                (transaction_id, 'income', 'income'),
                (transaction_id, 'salary', 'income');
        END IF;
        
        -- Healthcare transactions
        IF transaction_content ~* '(cvs|walgreens|pharmacy)' THEN
            INSERT INTO transactions_tags (transaction_id, tag_id, rule_id) VALUES 
                (transaction_id, 'healthcare', 'healthcare');
        END IF;
        
    END LOOP;
    
    RAISE NOTICE 'Successfully inserted 100 sample transactions with tags between % and %', year_start, current_date_val;
    
    -- Display summary
    RAISE NOTICE 'Transaction summary:';
    RAISE NOTICE '- Total transactions: %', (SELECT COUNT(*) FROM transactions WHERE date >= year_start);
    RAISE NOTICE '- Tagged transactions: %', (SELECT COUNT(DISTINCT transaction_id) FROM transactions_tags);
    RAISE NOTICE '- Total tag associations: %', (SELECT COUNT(*) FROM transactions_tags);
END;
$$ LANGUAGE plpgsql;

-- Execute the function to insert sample data
-- Usage: 
-- psql -h localhost -U postgres -d postgres -f tools/sample_transactions.sql
-- psql -h localhost -U postgres -d postgres -c "SELECT insert_sample_transactions();" 