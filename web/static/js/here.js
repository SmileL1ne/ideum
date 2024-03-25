var images = document.querySelectorAll('img.here');

images.forEach(function (image) {
    image.onerror = function () {
        // this.onerror=null
        this.src = '/static/img/svg/sport-icon.svg';
    };
});