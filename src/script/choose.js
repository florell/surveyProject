function showModal(title, description) {
    $("#modalLabel").text(title);
    $("#modalBody").html(description);
    $("#modal").modal("show");
    event.stopPropagation(); // предотвращает распространение события на родительский элемент
}
function redirectToSurvey(id) {
    window.location.href = "/survey/" + id;
}
function downloadTable() {
    $.ajax({
        url: '/get_table',
        type: 'GET',
        success: function () {
            // Создаем скрытую ссылку для скачивания файла
            const link = document.createElement('a');
            link.href = '/download_table';
            link.download = 'survey_results.xlsx';
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        },
        error: function () {
            alert('Ошибка при создании таблицы.');
        }
    });
}

function relogin(url) {
    window.location.href = "/";
}