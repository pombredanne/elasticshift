-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- +migrate StatementBegin
CREATE TABLE TEAMS (
    PUUID VARCHAR(32) PRIMARY KEY COMMENT 'UUID of TEAM',
    DOMAIN VARCHAR(63) NOT NULL COMMENT 'Subdomain of elasticshift.com', 
    NAME VARCHAR(255) NOT NULL COMMENT 'Displayable TEAM name',
    CREATED_DT DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'TEAM creation datetime',
    UPDATED_DT DATETIME ON UPDATE CURRENT_TIMESTAMP COMMENT 'TEAM updated datetime',
    CREATED_BY VARCHAR(100) COMMENT 'who created this record',
    UPDATED_BY VARCHAR(100) COMMENT 'who updated this record recently'
) COMMENT 'Subdomain or TEAM or team information';
-- +migrate StatementEnd


-------------------------------------------------------------------------------------------------------------------------------------------------
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE TEAMS;