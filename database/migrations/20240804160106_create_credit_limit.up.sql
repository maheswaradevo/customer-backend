CREATE TABLE IF NOT EXISTS credit_limits(
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,
    credit_limit DECIMAL(16, 2) NOT NULL,
    tenor_id INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (tenor_id) REFERENCES tenors(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);