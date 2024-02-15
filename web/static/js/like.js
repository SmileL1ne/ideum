document.addEventListener('DOMContentLoaded', function () {
    var likeButtons = document.querySelectorAll('.like-button');

    likeButtons.forEach(function (button) {
        button.addEventListener('click', function (event) {
            event.preventDefault(); // Отключаем отправку формы

            var modal = document.getElementById('signin');
            modal.showModal(); // Открываем модальное окно
        });
    });
});

