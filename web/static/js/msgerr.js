let LoginForm = document.querySelector('loginQuery')

form.event.preventDefault()

function LoginErr() {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/user/login", true);
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            // Обработка ответа от сервера
            // Можно выполнить какие-то действия после успешного лайка комментария
            console.log('OK')
        } else {
            console.log('NO')
        }
    };
    xhr.send();
}


document.addEventListener('DOMContentLoaded', function() {
    var form = document.querySelector('loginQuery');
    form.addEventListener('submit', function(event) {
        event.preventDefault();
        LoginErr();
    });
});
