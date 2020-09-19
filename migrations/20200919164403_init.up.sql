BEGIN;
CREATE TABLE received_messages (
	telegram_message_id int8 NOT NULL,
	user_id             int8 NOT NULL,
	chat_id             text NOT NULL,
	message_text        text NOT NULL
);
END;