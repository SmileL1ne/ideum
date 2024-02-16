function LoginErr(form) {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/user/login", true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                
                console.log('OK');
                window.location.href = "/"; 
            } else {
                console.log(xhr.status);
                var errorMsg = form.querySelector('.error-msg');
                if (errorMsg) {
                    errorMsg.textContent = 'Invalid Data';
                }
                
            }
        }
    };
    var formData = new FormData(form);
    var encodedData = new URLSearchParams(formData).toString(); 
    xhr.send(encodedData);
}

document.addEventListener('DOMContentLoaded', function() {
    var loginForms = document.querySelectorAll('.loginQuery');
    loginForms.forEach(function(form) {
        form.addEventListener('submit', function(event) {
            event.preventDefault();
            LoginErr(form);
        });
    });
});


function RegErr(form) {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/user/signup", true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                
                console.log('OK');
                window.location.href = "/"; 
            } else {
                console.log('NO');
                var errorMsg = form.querySelector('.error-msg');
                if (errorMsg) {
                    errorMsg.textContent = 'Invalid Data';
                }
                
            }
        }
    };
    var formData = new FormData(form); 
    var encodedData = new URLSearchParams(formData).toString(); 
    xhr.send(encodedData);
}

document.addEventListener('DOMContentLoaded', function() {
    var loginForms = document.querySelectorAll('.singupQuery');
    loginForms.forEach(function(form) {
        form.addEventListener('submit', function(event) {
            event.preventDefault();
            RegErr(form);
        });
    });
});
