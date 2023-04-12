CREATE TABLE workouts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    steps INT NOT NULL,
    calories INT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL
);
CREATE TABLE daily_aggregations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    steps INT NOT NULL,
    calories INT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL
);
CREATE TABLE weekly_aggregations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    steps INT NOT NULL,
    calories INT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL
);