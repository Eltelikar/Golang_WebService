CREATE TABLE IF NOT EXISTS subscriptions(
        id SERIAL PRIMARY KEY,
		service_name TEXT NOT NULL,
		price BIGINT CHECK (price > 0),
		user_id UUID NOT NULL,
		start_date DATE NOT NULL,
		end_date DATE,
		UNIQUE (service_name, user_id)
);

-- При использовании UNIQUE в запросе (8 строка), автоматически создается индекс на эти поля
-- Создание отдельного индекса может быть избыточным

CREATE INDEX IF NOT EXISTS idx_subscriptions_service_user ON subscriptions(service_name, user_id);
