package db

const (
	createUsersSQL = `
		CREATE TABLE IF NOT EXISTS users (
			user_id serial PRIMARY KEY,
			login varchar(64) UNIQUE NOT NULL,
			password varchar(64) NOT NULL
		)
	`
	createBalancesSQL = `
		CREATE TABLE IF NOT EXISTS balances (
			user_id int PRIMARY KEY REFERENCES users(user_id) NOT NULL,
			current	numeric(10,2) NOT NULL,
			withdrawn numeric(10,2) NOT NULL
		)
	`
	createOrdersSQL = `
		DO $$ BEGIN
	   		 CREATE TYPE status_enum AS ENUM ('NEW', 'REGISTERED', 'PROCESSING', 'INVALID', 'PROCESSED');
		EXCEPTION
    		WHEN duplicate_object THEN null;
		END $$;
		CREATE TABLE IF NOT EXISTS orders (
			number text UNIQUE NOT NULL,
			user_id int REFERENCES users(user_id) NOT NULL,
			status status_enum NOT NULL, 
			accrual	numeric(10,2),
			uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL
		)
	`
	createWithdrawalSQL = `
		CREATE TABLE IF NOT EXISTS withdrawal (
			number text UNIQUE NOT NULL,
			user_id int REFERENCES users(user_id) NOT NULL,
			withdrawn numeric(10,2),
			uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL
		)
	`
	getUserSQL = `
		SELECT user_id FROM users
		WHERE login = $1 AND password = $2
		LIMIT(1)
	`
	getUserForOrderSQL = `
		SELECT user_id FROM orders
		WHERE number = $1
		LIMIT(1)
	`
	getBalanceSQL = `
		SELECT current, withdrawn FROM balances
		WHERE user_id = $1
		LIMIT(1)
	`
	getCurrentBalanceSQL = `
	SELECT current FROM balances
	WHERE user_id = $1
	LIMIT(1)
`
	getOrdersSQL = `
		SELECT number, status, accrual, uploaded_at FROM orders
		WHERE user_id = $1
		ORDER BY uploaded_at DESC
	`
	getWithdrawalsSQL = `
		SELECT number, withdrawn, uploaded_at FROM withdrawal
		WHERE user_id = $1
		ORDER BY uploaded_at DESC
	`
	insertUserSQL = `
		INSERT INTO users (login, password)
		VALUES ($1, $2)
	`
	insertBalanceSQL = `
		INSERT INTO balances (user_id, current, withdrawn)
		VALUES ($1, $2, $3)
	`
	updateBalanceSQL = `
		UPDATE balances SET current = $2, withdrawn = $3 
		WHERE user_id = $1
	`
	updateCurrentBalanceSQL = `
		UPDATE balances SET current = $2
		WHERE user_id = $1
	`
	insertOrderSQL = `
			INSERT INTO orders (user_id, number, status, uploaded_at)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (number)
			DO UPDATE SET number = $2
				WHERE orders.user_id <> $1
			RETURNING user_id
		`

	/*
		insertOrderSQL = `
			INSERT INTO orders (user_id, number, status, uploaded_at)
			VALUES ($1, $2, $3, $4)
		`
	*/
	insertWithdrawalSQL = `
		INSERT INTO withdrawal (user_id, number, withdrawn, uploaded_at)
		VALUES ($1, $2, $3, $4)
	`
	updateOrderAccrualSQL = `
		UPDATE orders SET status = $2, accrual = $3 
		WHERE number = $1 AND status <> 'INVALID' AND status <> 'PROCESSED'
	`
	updateOrderStatusSQL = `
		UPDATE orders SET status = $2 
		WHERE number = $1 AND status <> 'INVALID' AND status <> 'PROCESSED'
	`
)
