CREATE TABLE signers (
    employee_id VARCHAR(10) PRIMARY KEY,
    username VARCHAR(255),
    firstname VARCHAR(255),
    lastname VARCHAR(255),
    campus VARCHAR(20),
    headers TEXT,
    time DATETIME
);
CREATE INDEX signers_username ON signers(username);
CREATE INDEX signers_firstname ON signers(firstname);
CREATE INDEX signers_lastname ON signers(lastname);
CREATE INDEX signers_campus ON signers(campus);
CREATE INDEX signers_time ON signers(time);
