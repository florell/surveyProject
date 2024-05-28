function saveValue(button) {
    let questionID = button.getAttribute('question');
    let value = button.getAttribute('value');
    sessionStorage.setItem('question_' + questionID, value);
}

let nextQuestion;
let currentQuestion = 0;

document.addEventListener("DOMContentLoaded", function() {
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

    nextQuestion = function() { 
        currentQuestion++;
        showQuestion(currentQuestion);
    };

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


    document.addEventListener('keydown', function(event) {
        if (event.key >= '1' && event.key <= '9') {
            var answerIndex = parseInt(event.key);
            var visibleQuestion = questions[currentQuestion];
            var answerButton = visibleQuestion.querySelector('.answer[data-answer-index="' + answerIndex + '"]');
            if (answerButton) {
                answerButton.click();
            }
        } else if (event.keyCode === 27) {
            let prevButton = document.getElementById('prev-button');
            prevButton.click();
        } else if (event.keyCode === 13) {
            let nextButton = document.getElementById('next-button');
            nextButton.click();
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

function saveValueFromFields() {
    let inputFields = document.querySelectorAll('.input-answer');
    let isValid = true;
    // убираем предыдущие уведомления перед переходом
    document.querySelectorAll('.error-message').forEach(function(errorDiv) {
        errorDiv.textContent = '';
    });
    for (let i = 0; i < inputFields.length; i++) {
        let field = inputFields[i];
        let value = parseInt(field.value);
        let min = parseInt(field.getAttribute('min'));
        let max = parseInt(field.getAttribute('max'));
        let errorDiv = document.getElementById('error-' + field.getAttribute('question'));

        if (!field.value.trim().length && i === currentQuestion) {
            console.log("Поле пустое")
            errorDiv.textContent = 'Поле не должно быть пустым.';
            isValid = false;
            break;
        }

        if (!isNaN(min) && !isNaN(max)) {
            if (value < min || value > max) {
                errorDiv.textContent = 'Значение должно быть в диапазоне от ' + min + ' до ' + max + '.';
                isValid = false;
                break;
            }
        } else if (!isNaN(min)) {
            if (value < min) {
                errorDiv.textContent = 'Значение должно быть не менее ' + min + '.';
                isValid = false;
                break; 
            }
        }

        sessionStorage.setItem('question_' + field.getAttribute('question'), field.value);
    }

    if (isValid) { 
        nextQuestion();
    }
}
