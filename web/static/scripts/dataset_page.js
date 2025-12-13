const inputSelectId = document.querySelector(".settings__input__select#id");
const inputSelectType = document.querySelector(".settings__input__select#type");
const inputDatetimeStart = document.querySelector(".settings__input__datetime#start");
const inputDatetimeEnd = document.querySelector(".settings__input__datetime#end");

const buildDatasetButton = document.querySelector(".buttons__dataset-button#build")
const downloadDatasetButton = document.querySelector(".buttons__dataset-button#download")

buildDatasetButton.addEventListener('click', async () => {
    const ids = Array.from(inputSelectId.selectedOptions).map(opt => opt.value).slice(1).toString();
    const types = Array.from(inputSelectType.selectedOptions).map(opt => opt.value).slice(1).toString();

    let startDate = new Date(inputDatetimeStart.value).getTime();
    let endDate = new Date(inputDatetimeEnd.value).getTime();

    if (endDate <= startDate) {alert("Конечная дата должна быть больше начальной!"); return;}
    if (isNaN(endDate)) {endDate = new Date().getTime();}
    if (isNaN(startDate)) {startDate = new Date().getTime() - 24*60*60*1000;}

    const url = `/api/get_dataset?startDate=${encodeURIComponent(startDate)}&endDate=${encodeURIComponent(endDate)}&id=${encodeURIComponent(ids)}&type=${encodeURIComponent(types)}`;

    const data = await request(url);
    console.log(data);
});

document.addEventListener('DOMContentLoaded', async function () {
});
