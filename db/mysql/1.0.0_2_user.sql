-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- +migrate StatementBegin
CREATE TABLE USER (
    ID VARCHAR(32) PRIMARY KEY COMMENT 'UUID of user',
    TEAM_ID VARCHAR(32) COMMENT 'UUID of TEAM',
    FULLNAME VARCHAR(70) NOT NULL COMMENT 'Full name of the user',
    USERNAME VARCHAR(100) NOT NULL COMMENT 'Username of shift-id',
    EMAIL VARCHAR(255) NOT NULL COMMENT 'Email of the user',
    PASSWORD VARCHAR(128) NOT NULL COMMENT 'Hashed password',
    LOCKED TINYINT DEFAULT 0 COMMENT '1-locked,0-unlocked',
    ACTIVE TINYINT DEFAULT 1 COMMENT '1-active,0-inactive',
    BAD_ATTEMPT TINYINT COMMENT 'Bad login attempt count',
    CREATED_DT DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'User creation datetime',
    UPDATED_DT DATETIME ON UPDATE CURRENT_TIMESTAMP COMMENT 'User updated datetime',
    CREATED_BY VARCHAR(100) COMMENT 'who created this record',
    UPDATED_BY VARCHAR(100) COMMENT 'who updated this record recently'
) COMMENT 'elasticshift user information';
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE USER
ADD CONSTRAINT FK_USER_TEAM_ID FOREIGN KEY (TEAM_ID)
REFERENCES TEAM(ID);
-- +migrate StatementEnd

-------------------------------------------------------------------------------------------------------------------------------------------------
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE USER DROP CONSTRAINT FK_USER_TEAM_ID;
DROP TABLE USER;