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
                const formattedText = xhr.responseText.replace(/\n/g, "<br><br>");
                errorMsg.innerHTML = formattedText;
            }
        };

        
        var encodedFormData = new URLSearchParams(formData).toString();

        xhr.send(encodedFormData);
    });
});
