CREATE TABLE breeds_characteristics_relation (
    breed_id INT NOT NULL,
    characteristic_id INT NOT NULL,
    PRIMARY KEY (breed_id, characteristic_id),
    FOREIGN KEY (breed_id) REFERENCES breeds(id) ON DELETE CASCADE,
    FOREIGN KEY (characteristic_id) REFERENCES characteristics(id) ON DELETE CASCADE
);