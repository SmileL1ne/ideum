const modalButtons = document.querySelectorAll('.Btn');
const closeButton = document.querySelectorAll('.BtnC');

modalButtons.forEach(button => {
    button.addEventListener('click', () => {
        const modalId = button.getAttribute('data-modal');
        const modal = document.getElementById(modalId);
        const urlId = button.getAttribute('data-url-id');

        if (modal) {
            if (urlId) {
                const form = modal.querySelector('form');
                form.action += `${urlId}`;
                modal.showModal();
            } else {
                modal.showModal();
            }

        }
    });
});

closeButton.forEach(button => {
    button.addEventListener('click', () => {
        const closeModal = button.closest('dialog');
        if (closeModal) {
            closeModal.close();
        }
    });
});

document.querySelectorAll('dialog').forEach(modal => {
    modal.addEventListener('click', (e) => {
        if (e.target === modal) modal.close();
    });
});





// class="Btn" data-modal="signin"
//  dialog -->id=signin

