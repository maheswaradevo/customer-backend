CREATE TABLE tenors (
    id INT AUTO_INCREMENT PRIMARY KEY,
    tenor_description ENUM('1 month', '2 months', '3 months', '6 months') NOT NULL UNIQUE
);