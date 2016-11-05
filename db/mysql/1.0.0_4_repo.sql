-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- +migrate StatementBegin
CREATE TABLE REPO (
    ID VARCHAR(32) PRIMARY KEY COMMENT 'UUID of vcs repository',
    TEAM_ID VARCHAR(32) NOT NULL COMMENT 'UUID of TEAM',
    VCS_ID VARCHAR(32) NOT NULL COMMENT 'UUID of VCS',
    REPO_ID VARCHAR(32) NOT NULL COMMENT 'Provider repository id',
    NAME VARCHAR(100) NOT NULL COMMENT 'Name of the repository',
    LANGUAGE VARCHAR(100) NOT NULL COMMENT 'Language of the repository',
    PRIVATE VARCHAR(1) NOT NULL COMMENT 'Private - Y, Public - N',
    LINK VARCHAR(256) NOT NULL COMMENT 'Http URL of the project',
    DESCRIPTION VARCHAR(512) NOT NULL COMMENT 'Repository Description',
    FORK VARCHAR(1) NOT NULL COMMENT 'Forked - Y, Othewise - N',
	DEFAULT_BRANCH VARCHAR(100) COMMENT 'Default branch of repository',
    CREATED_DT DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'Repository creation datetime',
    UPDATED_DT DATETIME ON UPDATE CURRENT_TIMESTAMP COMMENT 'Repository updated datetime',
    CREATED_BY VARCHAR(100) COMMENT 'who created this record',
    UPDATED_BY VARCHAR(100) COMMENT 'who updated this record recently'
) COMMENT 'Repository information';
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE REPO
ADD CONSTRAINT FK_REPO_TEAM_ID FOREIGN KEY (TEAM_ID)
REFERENCES TEAM(ID);
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE REPO
ADD CONSTRAINT FK_REPO_VCS_ID FOREIGN KEY (VCS_ID)
REFERENCES VCS(ID);
-- +migrate StatementEnd

-------------------------------------------------------------------------------------------------------------------------------------------------
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE REPO DROP CONSTRAINT FK_REPO_TEAM_ID;
ALTER TABLE REPO DROP CONSTRAINT FK_REPO_VCS_ID;

DROP TABLE REPO;
