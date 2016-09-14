-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- +migrate StatementBegin
CREATE TABLE VCS_ACCOUNTS (
    PUUID VARCHAR(32) PRIMARY KEY COMMENT 'UUID of vcs account',
    TEAM_PUUID VARCHAR(32) COMMENT 'UUID of TEAM',
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
ALTER TABLE VCS_ACCOUNTS
ADD CONSTRAINT FK_VCS_ACCOUNT_TEAM_PUUID FOREIGN KEY (TEAM_PUUID)
REFERENCES TEAMS(PUUID)
-- +migrate StatementEnd

-------------------------------------------------------------------------------------------------------------------------------------------------
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE VCS_ACCOUNTS DROP CONSTRAINT FK_VCS_ACCOUNT_TEAM_PUUID;

DROP TABLE VCS_ACCOUNTS;