CREATE TABLE IF NOT EXISTS customers(
    id INT NOT NULL AUTO_INCREMENT,
    nik_number VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255) NOT NULL,
    birthday_loc VARCHAR(100) NOT NULL,
    birthday_date DATE NOT NULL,
    salary DECIMAL(16, 2) NOT NULL,
    id_pic TEXT NULL,
    self_pic TEXT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY(id)
);