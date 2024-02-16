document.addEventListener('DOMContentLoaded', function () {
    var form = document.getElementById('postForm');

    form.addEventListener('submit', function (event) {
        event.preventDefault();

        var formData = new FormData(form);

        var xhr = new XMLHttpRequest();
        xhr.open('POST', form.action);

        
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

        
        var encodedFormData = new URLSearchParams(formData).toString();

        xhr.send(encodedFormData);
    });
});
