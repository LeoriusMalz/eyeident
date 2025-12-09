const mainLogo = document.querySelector(".page-header-navigation__logo img");

function updateNavigation() {
    const path = window.location.pathname.split('/')[1];
    const activeLink = document.querySelector(`.page-header-navigation__menu__link#${path}`);

    activeLink.classList.add("active");
    activeLink.disabled = true;
}

mainLogo.addEventListener('click', () => {
    window.location.href = '/';
});

document.addEventListener('DOMContentLoaded', function () {
    updateNavigation();
});