const inputSelectId = document.querySelector(".settings__input__select#id");
const inputSelectType = document.querySelector(".settings__input__select#type");
const inputDatetimeStart = document.querySelector(".settings__input__datetime#start");
const inputDatetimeEnd = document.querySelector(".settings__input__datetime#end");
const inputLimit = document.querySelector(".settings__input__limit")

const buildDatasetButton = document.querySelector(".buttons__dataset-button#build");
const downloadDatasetButton = document.querySelector(".buttons__dataset-button#download");

const loadingLayout = document.querySelector(".dataset-loading-layout");
const datasetAmount = document.querySelector(".preview__amount");
const datasetLabel = document.querySelector(".preview__label");
const table = document.querySelector(".preview__table");
const tableBody = document.querySelector(".preview__table__body");

let IDS, TYPES;

async function getParams() {
    const data = await request("/api/get_params");

    data['id'].forEach(params => {
        const option = document.createElement("option");
        option.textContent = params;

        inputSelectId.appendChild(option);
    });

    data['type'].forEach(params => {
        const option = document.createElement("option");
        option.textContent = params;

        inputSelectType.appendChild(option);
    });

    loadingLayout.style.display = 'none';
    return data;
}

function createTable(data) {
    datasetAmount.textContent = "Найдено "+data['count'].split(" ").slice(-1)+" записей";

    data['dataset'].forEach(row => {
        const tableRow = document.createElement("tr");

        const Id = document.createElement("td");
        Id.textContent = row["id"];

        const Ts = document.createElement("td");
        Ts.textContent = row["timestamp"];

        const Type = document.createElement("td");
        Type.textContent = row["type"];

        const AccelX = document.createElement("td");
        AccelX.textContent = row["accel_x"];

        const AccelY = document.createElement("td");
        AccelY.textContent = row["accel_y"];

        const AccelZ = document.createElement("td");
        AccelZ.textContent = row["accel_z"];

        const GyroX = document.createElement("td");
        GyroX.textContent = row["gyro_x"];

        const GyroY = document.createElement("td");
        GyroY.textContent = row["gyro_y"];

        const GyroZ = document.createElement("td");
        GyroZ.textContent = row["gyro_z"];

        const QX = document.createElement("td");
        QX.textContent = row["qx"];

        const QY = document.createElement("td");
        QY.textContent = row["qy"];

        const QZ = document.createElement("td");
        QZ.textContent = row["qz"];

        const QW = document.createElement("td");
        QW.textContent = row["qw"];

        const Yaw = document.createElement("td");
        Yaw.textContent = row["yaw"];

        const Pitch = document.createElement("td");
        Pitch.textContent = row["pitch"];

        const Roll = document.createElement("td");
        Roll.textContent = row["roll"];

        tableRow.appendChild(Id);
        tableRow.appendChild(Ts);
        tableRow.appendChild(Type);
        tableRow.appendChild(AccelX);
        tableRow.appendChild(AccelY);
        tableRow.appendChild(AccelZ);
        tableRow.appendChild(GyroX);
        tableRow.appendChild(GyroY);
        tableRow.appendChild(GyroZ);
        tableRow.appendChild(QX);
        tableRow.appendChild(QY);
        tableRow.appendChild(QZ);
        tableRow.appendChild(QW);
        tableRow.appendChild(Yaw);
        tableRow.appendChild(Pitch);
        tableRow.appendChild(Roll);

        tableBody.appendChild(tableRow);
    });
}

buildDatasetButton.addEventListener('click', async () => {
    let ids = Array.from(inputSelectId.selectedOptions).map(opt => opt.value).slice(1).toString();
    let types = Array.from(inputSelectType.selectedOptions).map(opt => opt.value).slice(1).toString();

    let startDate = new Date(inputDatetimeStart.value).getTime();
    let endDate = new Date(inputDatetimeEnd.value).getTime();

    let limit = inputLimit.value;

    if (ids === null || ids === "") {ids = IDS.toString();}
    if (types === null || types === "") {types = TYPES.toString();}

    if (endDate <= startDate) {alert("Конечная дата должна быть больше начальной!"); return;}
    if (isNaN(endDate)) {endDate = new Date().getTime();}
    if (isNaN(startDate)) {startDate = new Date().getTime() - 24*60*60*1000;}

    datasetLabel.textContent = "Сборка датасета...";
    datasetLabel.style.display = 'flex';

    const url = `/api/get_dataset?startDate=${encodeURIComponent(startDate)}&endDate=${encodeURIComponent(endDate)}&id=${encodeURIComponent(ids)}&type=${encodeURIComponent(types)}&limit=${limit}`;

    const data = await request(url);

    datasetLabel.style.display = 'none';
    tableBody.replaceChildren();
    createTable(data);
    table.style.display = 'block';

    datasetLabel.textContent = "Предварительный просмотр";
});

downloadDatasetButton.addEventListener('click', async () => {
    window.location.href = "/api/download_dataset";
    // await request("api/download_dataset");
})

document.addEventListener('DOMContentLoaded', async function () {
    const params = await getParams();
    IDS = params['id'];
    TYPES = params['type'];
});
