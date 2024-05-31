-- Создание функции для уведомлений
CREATE OR REPLACE FUNCTION notify_new_record() RETURNS trigger AS $$ DECLARE BEGIN PERFORM pg_notify('new_record', NEW::text); RETURN NEW; END; $$ LANGUAGE plpgsql;

-- Создание триггера для таблицы (предположим, таблица называется `items`)
CREATE TRIGGER new_record_trigger AFTER INSERT ON items FOR EACH ROW EXECUTE FUNCTION notify_new_record();
