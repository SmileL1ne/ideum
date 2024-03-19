console.log("here")
var images = document.querySelectorAll('img.here');

images.forEach(function (image) {
    image.onerror = function () {
        this.src = '/static/img/svg/sport-icon.svg';
    };
});