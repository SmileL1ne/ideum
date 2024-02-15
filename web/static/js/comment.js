document.addEventListener('DOMContentLoaded', function () {
    var form = document.getElementById('commentForm');

    form.addEventListener('submit', function (event) {
        event.preventDefault();

        var formData = new FormData(form);

        var xhr = new XMLHttpRequest();
        xhr.open('POST', form.action);

        // Устанавливаем заголовок Content-Type
        xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');

        xhr.onload = function () {
            if (xhr.status === 200) {
                console.log(xhr.responseText)
                window.location.href = xhr.responseText;
            } else {
                console.error('Request failed. Status: ' + xhr.status);
                var errorMsg = form.querySelector('.error-msg');
                if (errorMsg) {
                    errorMsg.textContent = 'Invalid Data';
                }
            }
        };

        // Преобразуем объект FormData в строку в формате application/x-www-form-urlencoded
        var encodedFormData = new URLSearchParams(formData).toString();

        xhr.send(encodedFormData);
    });
});
