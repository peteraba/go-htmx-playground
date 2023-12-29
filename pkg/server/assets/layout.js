// SSE - Server Sent Events
$(function () {
    const messages = $("#user-messages");

    const source = new EventSource(baseUrl + "/messages");
    source.onopen = (event) => {
        console.debug("/messages sse connection open");
    }
    source.onerror = (event) => {
        console.debug("/messages sse connection error");
    }

    const alert = (type, msg) => {
        const alert = $(`<p class="alert alert-${type}"><span>${msg}</span></p>`);
        alert.appendTo(messages);

        setTimeout(() => {alert.detach()}, 5000);
    }

    let firstConnection = true;
    const reload = () => {
        console.debug("reload received", window.liveReload);
        if (firstConnection) {
            firstConnection = false;
            console.debug("first connection, reloading skipped");
            return;
        }
        if (window.liveReload) {
            console.debug("reload to be executed", window.liveReload);
            window.location.reload();
        }
    }

    source.onmessage = (event) => {
        console.log("message received", event.data);
        const payload = JSON.parse(event.data),
            type = payload.type,
            msg = payload.message;

        switch (type) {
            case "reload":
                reload();
                break;
            default:
                alert(type, msg);
        }
    }
});
