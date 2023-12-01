$(function () {
    const theme = localStorage.getItem("theme");

    $("#html").attr("data-theme", theme);
});

$(function () {
    const messages = $("#user-messages");

    const source = new EventSource(baseUrl + "/sse");
    source.onopen = (event) => {
        console.log("sse connection open");
    }
    source.onerror = (event) => {
        console.error("sse connection error");
    }
    source.onmessage = (event) => {
        const payload = JSON.parse(event.data),
            type = payload.type,
            msg = payload.message;

        const alert = $(`<p class="alert alert-${type}"><span>${msg}</span></p>`);
        alert.appendTo(messages);

        setTimeout(() => {alert.detach()}, 5000);
    }
});
