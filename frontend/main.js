import {Events} from "@wailsio/runtime";

const timeElement = document.getElementById('time');

Events.On('time', (time) => {
    timeElement.innerText = time.data;
});
