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
/** @type HTMLButtonElement */
let trace_button = document.getElementById('trace-button');
/** @type HTMLDivElement */
let progress_tracker = document.getElementById('progress-tracker');


trace_button.addEventListener(
    "click",
    async () => {
        for await (let chunk of streamingFetch(() => fetch("/api/stream"))) {
            parts = chunk.toString().split("/");
            done = parseInt(parts[0]);
            full = parseInt(parts[1]);
            precent = done / full * 100;
            progress_tracker.innerText = `Progress: ${precent.toFixed(2)}%`;
        }
    }
)

submit_button.addEventListener(
    "click",
    () => {
        fetch("/api/download", {
            method: "POST",
            body: JSON.stringify({ url: url_input.value }),
        })
    }
)
