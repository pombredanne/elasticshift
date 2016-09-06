-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- +migrate StatementBegin
CREATE TABLE SITE (
    PUUID VARCHAR(32) PRIMARY KEY COMMENT 'UUID of site',
    DOMAIN VARCHAR(63) NOT NULL COMMENT 'Subdomain of elasticshift.com', 
    NAME VARCHAR(255) NOT NULL COMMENT 'Displayable site name',
    CREATED_DT DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'Site creation datetime',
    UPDATED_DT DATETIME ON UPDATE CURRENT_TIMESTAMP COMMENT 'Site updated datetime',
    CREATED_BY VARCHAR(100) COMMENT 'who created this record',
    UPDATED_BY VARCHAR(100) COMMENT 'who updated this record recently'
) COMMENT 'Subdomain or site or team information';
-- +migrate StatementEnd


-- +migrate StatementBegin
CREATE TABLE USER (
    PUUID VARCHAR(32) PRIMARY KEY COMMENT 'UUID of user',
    SITE_PUUID VARCHAR(32) COMMENT 'UUID of site',
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
ALTER TABLE USER
ADD CONSTRAINT FK_USER_SITE_PUUID FOREIGN KEY (SITE_PUUID)
REFERENCES SITE(PUUID);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE VCS_ACCOUNT (
    PUUID VARCHAR(32) PRIMARY KEY COMMENT 'UUID of vcs account',
    SITE_PUUID VARCHAR(32) COMMENT 'UUID of site',
    NAME VARCHAR(100) NOT NULL COMMENT 'username of version control system',
    TYPE TINYINT DEFAULT 0 COMMENT '1-github, 2-gitlab, 3-bitbucket, 4-SVN, 5-TFS',
    HASHED_CODE VARCHAR(255) NOT NULL COMMENT 'Hashed access code',
    REVOKED TINYINT DEFAULT 0 COMMENT '1-revoked, 0-active',
    CREATED_DT DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'VCS account creation datetime',
    UPDATED_DT DATETIME ON UPDATE CURRENT_TIMESTAMP COMMENT 'VCS account updated datetime',
    CREATED_BY VARCHAR(100) COMMENT 'who created this record',
    UPDATED_BY VARCHAR(100) COMMENT 'who updated this record recently'
) COMMENT 'Version control system information';
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE VCS_ACCOUNT
ADD CONSTRAINT FK_VCS_ACCOUNT_SITE_PUUID FOREIGN KEY (SITE_PUUID)
REFERENCES SITE(PUUID)
-- +migrate StatementEnd

-------------------------------------------------------------------------------------------------------------------------------------------------
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE VCS_ACCOUNT DROP CONSTRAINT FK_SITE_PUUID;
ALTER TABLE USER DROP CONSTRAINT FK_SITE_PUUID;

DROP TABLE VCS_ACCOUNT;
DROP TABLE USER;
DROP TABLE SITE;
