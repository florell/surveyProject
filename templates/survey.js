function saveValue(button) {
    let questionID = button.getAttribute('question');
    let value = button.getAttribute('value');
    sessionStorage.setItem('question_' + questionID, value);
}

function saveValueFromFields() {
    let inputFields = document.querySelectorAll('.input-answer');
    inputFields.forEach(function(field) {
        sessionStorage.setItem('question_' + field.getAttribute('question'), field.value);
    });
}

document.addEventListener("DOMContentLoaded", function() {
    let currentQuestion = 0;
    let questions = document.querySelectorAll('.question');
    let prevButton = document.getElementById('prev-button');
    let nextButton = document.getElementById('next-button');
    let questionCount = parseInt(document.getElementById('questionCount').value);
    let firstQuestionID = parseInt(document.getElementById('firstQuestionID').value);

    function showQuestion(questionIndex) {
        if (questionIndex === questionCount) {
            prevButton.style.display = 'none';
            submitForm();
        } else {
            for (let i = 0; i < questions.length; i++) {
                questions[i].style.display = 'none';
            }
            if (questionIndex === 0) {
                prevButton.style.display = 'none';
            } else {
                prevButton.style.display = 'block';
            }
            questions[questionIndex].style.display = 'block';
        }
    }

    function nextQuestion() {
        currentQuestion++;
        showQuestion(currentQuestion);
    }

    function prevQuestion() {
        if (currentQuestion > 0) {
            currentQuestion--;
            showQuestion(currentQuestion);
        }
    }

    showQuestion(currentQuestion);

    let answerButtons = document.querySelectorAll('.answer');
    answerButtons.forEach(function(button) {
        button.addEventListener('click', function() {
            nextQuestion();
        });
    });

    prevButton.addEventListener('click', function() {
        prevQuestion();
    });
    if (nextButton != null) {
        nextButton.addEventListener('click', function() {
            nextQuestion();
        });
    }

    document.addEventListener('keydown', function(event) {
        if (event.key >= '1' && event.key <= '9') {
            var answerIndex = parseInt(event.key);
            var visibleQuestion = questions[currentQuestion];
            var answerButton = visibleQuestion.querySelector('.answer[data-answer-index="' + answerIndex + '"]');
            if (answerButton) {
                answerButton.click();
            }
        }
        if (event.keyCode === 27) {
            let prevButton = document.getElementById('prev-button');
            prevButton.click();
        }
    });
});

function submitForm() {
    event.preventDefault();
    let form = document.getElementById('survey-form');
    for (let i = parseInt(firstQuestionID.value); i < parseInt(questionCount.value) + parseInt(firstQuestionID.value); i++) {
        let answerValue = sessionStorage.getItem('question_' + i);

        let hiddenInput = document.createElement('input');
        hiddenInput.type = 'hidden';
        hiddenInput.name = 'question_' + i;
        hiddenInput.value = answerValue;

        form.appendChild(hiddenInput);
    }

    form.submit();
}