-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- +migrate StatementBegin
CREATE TABLE VCS (
    ID VARCHAR(32) PRIMARY KEY COMMENT 'UUID of vcs account',
    TEAM_ID VARCHAR(32) COMMENT 'UUID of TEAM',
    NAME VARCHAR(100) NOT NULL COMMENT 'username of version control system',
    TYPE TINYINT DEFAULT 0 COMMENT '1-github, 2-gitlab, 3-bitbucket, 4-SVN, 5-TFS',
    OWNER_TYPE TINYINT DEFAULT 0 COMMENT '1-user, 2-org',
    ACCESS_TOKEN VARCHAR(255) NOT NULL COMMENT 'vcs access code',
    AVATAR_URL VARCHAR(255) NOT NULL COMMENT 'avatar url',
    REFRESH_TOKEN VARCHAR(255) NOT NULL COMMENT 'refresh token',
    TOKEN_TYPE VARCHAR(64) NOT NULL COMMENT 'token type',
    TOKEN_EXPIRY DATETIME  NULL COMMENT "Token expiry time",
    REVOKED TINYINT DEFAULT 0 COMMENT '1-revoked, 0-active',
    CREATED_DT DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'VCS account creation datetime',
    UPDATED_DT DATETIME ON UPDATE CURRENT_TIMESTAMP COMMENT 'VCS account updated datetime',
    CREATED_BY VARCHAR(100) COMMENT 'who created this record',
    UPDATED_BY VARCHAR(100) COMMENT 'who updated this record recently'
) COMMENT 'Version control system information';
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE VCS
ADD CONSTRAINT FK_VCS_TEAM_ID FOREIGN KEY (TEAM_ID)
REFERENCES TEAM(ID)
-- +migrate StatementEnd

-------------------------------------------------------------------------------------------------------------------------------------------------
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE VCS DROP CONSTRAINT FK_VCS_TEAM_ID;

DROP TABLE VCS;