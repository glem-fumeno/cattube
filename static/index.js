/** @typedef StreamedData
 *  @type {object}
 *  @property {string} currentLog
 *  @property {number} currentSize
 *  @property {number} currentTotalSize
 *  @property {SDNode[]} videos
 */

/** @typedef SDNode
 * @type {object}
 * @property {string} url
 * @property {string} title
 * @property {string} duration
 * */

async function* streamingFetch(fetchcall) {
    const response = await fetchcall();
    const reader = response.body.getReader();
    while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        yield (new TextDecoder().decode(value));
    }
}


/** @type HTMLInputElement */
let url_input = document.getElementById('input-url');
/** @type HTMLButtonElement */
let submit_button = document.getElementById('submit-button');
/** @type HTMLDivElement */
let progress_tracker = document.getElementById('progress-tracker');
/** @type HTMLElement */
const progress_tracker_tracker = document.querySelector("#progress-tracker-tracker");

const video_queue = document.querySelector("#video-queue")

submit_button.addEventListener(
    "click",
    async () => {
        document.querySelector("#server-response").innerText = "Downloading...";
        fetch("/api/download", {
            method: "POST",
            body: JSON.stringify({ url: url_input.value }),
        })
        setTimeout(getStreamAndParseResponse, 500);
    }
)

window.onload = getStreamAndParseResponse;

async function getStreamAndParseResponse() {
    for await (let chunk of streamingFetch(() => fetch("/api/stream"))) {
        /** @type {StreamedData} */
        const sd = JSON.parse(chunk);
        let log = sd.currentLog;
        const done = sd.currentSize;
        const full = sd.currentTotalSize;
        let percent = done / full * 100;
        if (isNaN(percent) || !isFinite(percent)) {
            percent = 0;
        }
        if (log.length > 100) {
            log = log.substring(0, 100) + "...";
        }
        document.querySelector("#server-response").innerText = log;
        document.querySelector("p").innerText = `${percent.toFixed(2)}%`;
        progress_tracker_tracker.style.width = parseInt(percent) + "%";

        video_queue.innerHTML = "";

        const videos = sd.videos;

        for (const video of videos) {
            const li = document.createElement("li");
            li.innerHTML = `${video.title}`;
            video_queue.appendChild(li);
        }
    }
}
