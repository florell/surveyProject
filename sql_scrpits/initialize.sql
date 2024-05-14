CREATE DATABASE IF NOT EXISTS psy_data;
USE psy_data;

-- Create Surveys table
CREATE TABLE IF NOT EXISTS surveys (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Title VARCHAR(255),
    Description VARCHAR(255)
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
    Surname VARCHAR(255),
    Age INT,
    Sex VARCHAR(255),
    UNIQUE KEY unique_columns (Name, Surname, Age, Sex)
);

CREATE TABLE IF NOT EXISTS survey_results (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    PatientID INT,
    SurveyID INT,
    Result JSON,
    CurDate Date,
    FOREIGN KEY (PatientID) REFERENCES patients(ID),
    FOREIGN KEY (SurveyID) REFERENCES surveys(ID),
    UNIQUE KEY unique_columns (PatientID, SurveyID)
);

