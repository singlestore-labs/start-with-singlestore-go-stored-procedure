CREATE DATABASE IF NOT EXISTS acme;

USE acme;

CREATE TABLE IF NOT EXISTS messages (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  content varchar(300) NOT NULL,
  createdate TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

DELIMITER //

CREATE PROCEDURE messages_create(_content varchar(300))
AS
BEGIN
  INSERT INTO messages (
    content
  ) VALUES (
    _content
  );
  ECHO SELECT LAST_INSERT_ID() as id;
END //

CREATE PROCEDURE messages_read_by_id(_id bigint)
AS
BEGIN
  ECHO SELECT id, content, createdate FROM messages WHERE id = _id;
END //

CREATE PROCEDURE messages_read_all()
AS
BEGIN
  ECHO SELECT * FROM messages ORDER BY id;
END //

CREATE PROCEDURE messages_update(_id bigint, _content varchar(300))
AS
BEGIN
  UPDATE messages SET content = _content WHERE id = _id;
END //

CREATE PROCEDURE messages_delete(_id bigint)
AS
BEGIN
  DELETE FROM messages WHERE id = _id;
END //

DELIMITER ;

CALL messages_create('starting docker container');
