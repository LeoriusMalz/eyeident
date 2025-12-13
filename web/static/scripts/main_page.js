const mainLogo = document.querySelector(".page-header-navigation__logo img");

async function request(url, options = {}) {
    try {
        const response = await fetch(url, options);

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`);
        }

        if (response.redirected) {
            window.location.href = response.url;
        }

        const data = await response.json();

        return data;
    } catch (error) {
        console.error("Request failed:", error);
        throw error;
    }
}

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