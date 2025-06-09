CREATE TABLE vacancies (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	title TEXT NOT NULL,
	company_name TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL
);

