-- Create Surveys table
CREATE TABLE IF NOT EXISTS surveys (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Title VARCHAR(255)
);

-- Create Questions table
CREATE TABLE IF NOT EXISTS questions (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    SurveyID INT,
    Title VARCHAR(255),
    FOREIGN KEY (SurveyID) REFERENCES surveys(ID)
);

-- Create Answers table
CREATE TABLE IF NOT EXISTS answers (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    QuestionID INT,
    Text VARCHAR(255),
    Value INT,
    FOREIGN KEY (QuestionID) REFERENCES questions(ID)
);

-- Create Patients table
CREATE TABLE IF NOT EXISTS patients (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Name VARCHAR(255),
    Surname VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS survey_results (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    PatientID INT,
    SurveyID INT,
    Result JSON,
    FOREIGN KEY (PatientID) REFERENCES patients(ID),
    FOREIGN KEY (SurveyID) REFERENCES surveys(ID)
);

