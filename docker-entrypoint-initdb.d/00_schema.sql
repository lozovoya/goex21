CREATE TABLE companies
(
    ID BIGSERIAL PRIMARY KEY,
    Name TEXT UNIQUE NOT NULL CHECK ( Name != '' ),
    Code TEXT UNIQUE NOT NULL CHECK ( Code != '' ),
    Country TEXT NOT NULL,
    Website TEXT NOT NULL,
    Phone TEXT NOT NULL,
    IsActive BOOL DEFAULT TRUE NOT NULL,
    Created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);