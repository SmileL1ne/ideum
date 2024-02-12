document.addEventListener('DOMContentLoaded', function () {
    // Получаем все элементы с классом "my-button"
    var buttons = document.querySelectorAll('.like-button');

    // Добавляем обработчик события для каждой кнопки
    buttons.forEach(function (button) {
        button.addEventListener('click', function () {
            // Сбросить состояние предыдущей кнопки
            buttons.forEach(function (otherButton) {
                if (otherButton !== button) {
                    otherButton.classList.remove('liked-button');
                }
            });

            // Ваш код для обработки нажатия текущей кнопки
            this.classList.toggle('liked-button'); // добавляем/удаляем класс "red-button"
        });
    });
});
