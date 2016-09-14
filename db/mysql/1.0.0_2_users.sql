-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- +migrate StatementBegin
CREATE TABLE USERS (
    PUUID VARCHAR(32) PRIMARY KEY COMMENT 'UUID of user',
    TEAM_PUUID VARCHAR(32) COMMENT 'UUID of TEAM',
    FIRSTNAME VARCHAR(35) NOT NULL COMMENT 'First name of the user',
    LASTNAME VARCHAR(35) NOT NULL COMMENT 'Last name of the user',
    USERNAME VARCHAR(100) NOT NULL COMMENT 'Username of shift-id',
    EMAIL VARCHAR(255) NOT NULL COMMENT 'Email of the user',
    HASHED_PASSWORD VARCHAR(128) NOT NULL COMMENT 'Hashed password',
    LOCKED TINYINT DEFAULT 0 COMMENT '1-locked,0-unlocked',
    ACTIVE TINYINT DEFAULT 1 COMMENT '1-active,0-inactive',
    BAD_ATTEMPT TINYINT COMMENT 'Bad login attempt count',
    LAST_LOGIN DATETIME COMMENT 'User last login datetime',
    VERIFY_CODE MEDIUMINT COMMENT 'verification code',
    CREATED_DT DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'User creation datetime',
    UPDATED_DT DATETIME ON UPDATE CURRENT_TIMESTAMP COMMENT 'User updated datetime',
    CREATED_BY VARCHAR(100) COMMENT 'who created this record',
    UPDATED_BY VARCHAR(100) COMMENT 'who updated this record recently'
) COMMENT 'elasticshift user information';
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE USERS
ADD CONSTRAINT FK_USER_TEAM_PUUID FOREIGN KEY (TEAM_PUUID)
REFERENCES TEAMS(PUUID);
-- +migrate StatementEnd

-------------------------------------------------------------------------------------------------------------------------------------------------
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE USERS DROP CONSTRAINT FK_USER_TEAM_PUUID;
DROP TABLE USERS;