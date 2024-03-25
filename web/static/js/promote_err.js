document.addEventListener('DOMContentLoaded', function () {
   
    var promoteLink = document.getElementById('promote');


    function sendRequest() {
        var xhr = new XMLHttpRequest();
        xhr.open('GET', promoteLink.href, true);
        xhr.onload = function () {
          
            var errorElement = document.getElementById('small-err');
           
            if (xhr.status === 200) {
            
                if (errorElement) {
                    errorElement.parentNode.removeChild(errorElement);
                }
            } else {
           
                if (!errorElement) {
                    errorElement = document.createElement('p');
                    errorElement.className = 'error-msg';
                    errorElement.id = 'small-err';
                    promoteLink.parentNode.insertBefore(errorElement, promoteLink);
                }
                 text = xhr.responseText
                errorElement.textContent = text;
            }
        };
        xhr.onerror = function () {
          
            var errorElement = document.getElementById('small-err');
            if (!errorElement) {
                errorElement = document.createElement('p');
                errorElement.className = 'error-msg';
                errorElement.id = 'small-err';
                promoteLink.parentNode.insertBefore(errorElement, promoteLink);
            }
            
            errorElement.textContent = 'Connection Error';
        };
        xhr.send();
    }

    
    promoteLink.addEventListener('click', function (event) {
        event.preventDefault(); 
        sendRequest(); 
    });
});