CREATE TABLE signers (
    username VARCHAR(255) PRIMARY KEY,
    firstname VARCHAR(255),
    lastname VARCHAR(255),
    campus VARCHAR(20),
    headers TEXT,
    time DATETIME
);
CREATE INDEX signers_firstname ON signers(firstname);
CREATE INDEX signers_lastname ON signers(lastname);
CREATE INDEX signers_campus ON signers(campus);
CREATE INDEX signers_time ON signers(time);
