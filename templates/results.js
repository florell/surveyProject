const surveyDataString = document.getElementById("surveyData").innerText;
const surveyDataObject = JSON.parse(surveyDataString);
console.log(surveyDataObject);

function renderProgressBars(data) {
    const featureProgressContainer = document.getElementById('feature-progress');
    const progressBarContainer = document.getElementById('progress-bars');

    for (const key in data) {
        if (data.hasOwnProperty(key)) {
            const feature = data[key];
            // Check if the key is 'code'
            if (key === 'code') {
                // Create a separate field for 'code'
                const codeHTML = `
                    <div class="progress-bar-container">
                        <div class="progress-bar-label">Код профиля:</div>
                        <div class="progress-bar">
                            <div class="progress-bar-inner" style="--hue: 120; width: 100%;">
                                ${feature}
                            </div>
                        </div>
                    </div>`;
                featureProgressContainer.innerHTML += codeHTML;
            } else {
                // For other keys, continue as before
                let progress = (feature.value / feature.max_value) * 100;
                if (progress < 10 || feature.value == null) {
                    progress = 10;
                }
                if (progress >= 100) {
                    progress = 92;
                }
                const hue = (progress / 100) * 120;

                let percent_string = "";
                if (feature.percent != null) {
                    percent_string = " (" + feature.percent.toFixed(2) + "%) "
                }

                let t_score_string = "";
                if (feature.tscore != null) {
                    t_score_string = " (" + feature.tscore + " Т) ";
                }

                let value = 0;
                if (feature.value != null) {
                    value = feature.value;
                }
                let result_string = value + percent_string + t_score_string;

                if (feature.value === null) {
                    feature.value = 0
                }

                const progressBarHTML = `
                    <div class="progress-bar-container">
                        <div class="progress-bar-label">${key}</div>
                        <div class="progress-bar">
                            <div class="progress-bar-gray"></div>
                            <div class="progress-bar-inner" style="--hue: ${hue}; width: ${progress}%;">
                                ${result_string}
                            </div>
                            <div class="progress-bar-max-value">${feature.max_value}</div>

                        </div>
                    </div>`;
                if (key === 'Feature progress') {
                    featureProgressContainer.innerHTML += progressBarHTML;
                } else {
                    progressBarContainer.innerHTML += progressBarHTML;
                }
            }
        }
    }
}

// Call the fetchData function when the page loads
window.onload = renderProgressBars(surveyDataObject);