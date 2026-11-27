DELETE FROM bank_exchange_rates;

INSERT
	INTO
	bank_exchange_rates (exchange_rate_uuid,
	from_currency,
	to_currency,
	rate,
	valid_from_timestamp,
	valid_to_timestamp,
	created_at,
	updated_at)
VALUES('31f9f0bd-86e7-46a8-918f-57242aa04f75',
'USD',
'UAH',
44.3452445,
'2026-11-20T15:54:24.648267',
'2026-11-28T15:54:24.648267',
now(),
now())
ON CONFLICT DO NOTHING;


