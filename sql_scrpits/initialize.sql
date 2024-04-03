-- Create Surveys table
CREATE TABLE surveys (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Title VARCHAR(255)
);

-- Create Questions table
CREATE TABLE questions (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    SurveyID INT,
    Title VARCHAR(255),
    FOREIGN KEY (SurveyID) REFERENCES surveys(ID)
);


-- Create Answers table
CREATE TABLE answers (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    QuestionID INT,
    Text VARCHAR(255),
    Value VARCHAR(255),
    FOREIGN KEY (QuestionID) REFERENCES questions(ID)
);
