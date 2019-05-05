CREATE TABLE Location (
    id       INTEGER        PRIMARY KEY AUTO_INCREMENT,
    name     VARCHAR(100)   NOT NULL,
    longlat VARCHAR(50)    NOT NULL UNIQUE
);
CREATE TABLE Reading (
    id       INTEGER        PRIMARY KEY AUTO_INCREMENT,
    stamp    VARCHAR(100)   NOT NULL UNIQUE,
    ts       TIMESTAMP  NOT NULL,
    windspeed INTEGER,
    temperature INTEGER,
    humidity INTEGER,
    yyyymmdd    DATE,
    location INTEGER,
    FOREIGN KEY (location) REFERENCES Location(id)
);

CREATE TABLE Forecast (
    id       INTEGER        PRIMARY KEY AUTO_INCREMENT,
    stamp    VARCHAR(100)   NOT NULL UNIQUE,
    ts       TIMESTAMP  NOT NULL,
    cond VARCHAR(50),
    current INTEGER,
    hightemp INTEGER,
    lowtemp INTEGER,
    humidity INTEGER,
    windspeed INTEGER,
    maxWindspeed INTEGER,
    winddirection VARCHAR(50),
    yyyymmdd    DATE,
    location INTEGER,
    FOREIGN KEY (location) REFERENCES Location(id)
);

CREATE TABLE User (
    id       INTEGER        PRIMARY KEY AUTO_INCREMENT,
    ts       TIMESTAMP  NOT NULL,
    email VARCHAR(50),
    username VARCHAR(100),
    frequency VARCHAR(50),
    content VARCHAR(100),
    location INTEGER,
    FOREIGN KEY (location) REFERENCES Location(id)
);

INSERT INTO Location (name, longlat) VALUES ( 'Bristol', '51.441,-2.57');
INSERT INTO Reading (stamp, windspeed, temperature, humidity, yyyymmdd) VALUES ( 'May 18, 1919 at 07:00AM', 9,12,67,"1919-05-18");
INSERT INTO Forecast (stamp, cond, current, hightemp, lowtemp, humidity, windspeed, maxwindspeed,winddirection,yyyymmdd,location) VALUES ( 'May 21, 2018 at 07:00AM',"Sweaty",12, 18,12,67,12,22,"Northwest","2018-05-22",1);
