const userListTable = document.querySelector(".users-list__table");

async function getUsers() {
    const data = await request("/api/get_users");
    data.forEach(user => {
        user['LastCommit'] = getLastCommit(user['LastCommit'])
        addUserCell(user);
    });

    const pauseButtons = document.querySelectorAll(".pause-button");
    pauseButtons.forEach(button => {
        button.addEventListener('click', async (e) => {
            const icon = button.querySelector("i");
            const options = {method: 'POST',};

            if (icon.classList.contains("fa-pause")) {
                icon.classList.remove("fa-pause");
                icon.classList.add("fa-play");
                button.title = 'Запустить сбор';

                await request(`/api/change_user_disable?id=${button.parentElement.parentElement.id}`, options);
            } else {
                icon.classList.remove("fa-play");
                icon.classList.add("fa-pause");
                button.title = 'Приостановить сбор';

                await request(`/api/change_user_enable?id=${button.parentElement.parentElement.id}`, options);
            }
        });
    });
}

function getLastCommit(lastCommit) {
    const lastCommitMs = Date.parse(lastCommit);
    let diffSec = parseInt(((Date.now() - lastCommitMs) / 1000).toFixed(0));
    let lastCommitTxt;

    if (diffSec < 60) {
        lastCommitTxt = `${diffSec} сек. назад`;
    } else if (diffSec < 3600) {
        diffSec = (diffSec / 60).toFixed(0);
        lastCommitTxt = `${diffSec} мин. назад`;
    } else if (diffSec < 86400) {
        diffSec = (diffSec / 3600).toFixed(0);
        lastCommitTxt = `${diffSec} ч. назад`;
    } else if (diffSec < 2592000) {
        diffSec = (diffSec / 86400).toFixed(0);
        lastCommitTxt = `${diffSec} д. назад`;
    } else {
        lastCommitTxt = "Больше месяца назад";
    }

    return [lastCommitTxt, new Date(lastCommitMs)];
}

function addUserCell(user_info) {
    const tableRow = document.createElement("tr");
    tableRow.classList.add("users-list__table__element");
    tableRow.id = user_info['Id'];

    let isActive, isActiveTitle, isEnabled, isEnabledTitle;
    if (user_info['IsActive']) {isActive = 'active'} else {isActive = 'inactive'}
    if (user_info['IsActive']) {isActiveTitle = 'Подключен'} else {isActiveTitle = 'Отключен'}
    if (user_info['IsEnabled']) {isEnabled = 'fa-pause'} else {isEnabled = 'fa-play'}
    if (user_info['IsEnabled']) {isEnabledTitle = 'Приостановить сбор'} else {isEnabledTitle = 'Запустить сбор'}


    const tableDataStatus = document.createElement("td");
    tableDataStatus.classList.add("element__status", isActive);

    const tableDataStatusCircle = document.createElement("i");
    tableDataStatusCircle.classList.add("fa", "fa-circle");
    tableDataStatusCircle.title = isActiveTitle;

    tableDataStatus.appendChild(tableDataStatusCircle);


    const tableDataUserId = document.createElement("td");
    tableDataUserId.classList.add("element__user-id");
    tableDataUserId.textContent = user_info['Id'];


    const tableDataLastCommit = document.createElement("td");
    tableDataLastCommit.classList.add("element__last-commit");
    tableDataLastCommit.title = user_info['LastCommit'][1].toString();
    tableDataLastCommit.textContent = user_info['LastCommit'][0];


    const tableDataLastStatus = document.createElement("td");
    tableDataLastStatus.classList.add("element__last-status");
    tableDataLastStatus.title = user_info['LastStatus'];
    tableDataLastStatus.textContent = user_info['LastStatus'].slice(0, 20);
    if (user_info['LastStatus'].length > 20) {tableDataLastStatus.textContent += '...'}


    const tableDataLogButton = document.createElement("td");
    tableDataLogButton.classList.add("element__logs");

    const logButton = document.createElement("button");
    logButton.classList.add("logs-button");
    logButton.title = "Смотреть логи";
    logButton.textContent = "Логи";

    tableDataLogButton.appendChild(logButton);


    const tableDataDownloadButton = document.createElement("td");
    tableDataDownloadButton.classList.add("element__download-button");

    const downloadButton = document.createElement("button");
    downloadButton.classList.add("download-button");
    downloadButton.title = "Скачать последние 24 часов";

    const downloadButtonSign = document.createElement("i");
    downloadButtonSign.classList.add("fa", "fa-download");

    downloadButton.appendChild(downloadButtonSign);
    tableDataDownloadButton.appendChild(downloadButton);


    const tableDataCounterEnableButton = document.createElement("td");
    tableDataCounterEnableButton.classList.add("element__counter-enable-button");

    const counterEnableButton = document.createElement("button");
    counterEnableButton.classList.add("pause-button");
    counterEnableButton.title = isEnabledTitle;

    const counterEnableButtonSign = document.createElement("i");
    counterEnableButtonSign.classList.add("fa", isEnabled);

    counterEnableButton.appendChild(counterEnableButtonSign);
    tableDataCounterEnableButton.appendChild(counterEnableButton);


    tableRow.appendChild(tableDataStatus);
    tableRow.appendChild(tableDataUserId);
    tableRow.appendChild(tableDataLastCommit);
    tableRow.appendChild(tableDataLastStatus);
    tableRow.appendChild(tableDataLogButton);
    tableRow.appendChild(tableDataDownloadButton);
    tableRow.appendChild(tableDataCounterEnableButton);

    userListTable.appendChild(tableRow);
}

document.addEventListener('DOMContentLoaded', async function () {
    await getUsers();
});