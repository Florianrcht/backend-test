CREATE TABLE breeds (
    id INT AUTO_INCREMENT PRIMARY KEY,
    species ENUM('dog', 'cat') NOT NULL,
    pet_size ENUM('small', 'medium', 'tall') NOT NULL,
    name VARCHAR(50) NOT NULL,
    average_male_adult_weight INT NOT NULL,
    average_female_adult_weight INT NOT NULL,
    hypoallergenic BOOLEAN NOT NULL,
    apartment_friendly BOOLEAN NOT NULL
);