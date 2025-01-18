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


submit_button.addEventListener(
    "click",
    async () => {
        document.querySelector("#server-response").innerText = "Downloading...";
        fetch("/api/download", {
            method: "POST",
            body: JSON.stringify({ url: url_input.value }),
        }).then(
            async (response) => {
                const data = await response.text();
                document.querySelector("#server-response").innerText = data;
            }
        )
        for await (let chunk of streamingFetch(() => fetch("/api/stream"))) {
            parts = chunk.toString().split("/");
            done = parseInt(parts[0]);
            full = parseInt(parts[1]);
            precent = done / full * 100;
            document.querySelector("p").innerText = `${precent.toFixed(2)}%`;

            progress_tracker_tracker.style.width = parseInt(precent) + "%";
        }
    }
)

window.onload = async () => {
    for await (let chunk of streamingFetch(() => fetch("/api/stream"))) {
        parts = chunk.toString().split("/");
        done = parseInt(parts[0]);
        full = parseInt(parts[1]);
        precent = done / full * 100;
        document.querySelector("p").innerText = `${precent.toFixed(2)}%`;

        progress_tracker_tracker.style.width = parseInt(precent) + "%";
    }
}
