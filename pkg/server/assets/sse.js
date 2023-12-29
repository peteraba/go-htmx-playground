// SSE - Server Sent Events
(function ready(fn) {
    if (document.readyState !== 'loading') {
        fn();
    } else {
        document.addEventListener('DOMContentLoaded', fn);
    }
})(
    () => {
        const alertContainer = document.getElementById("user-messages");

        const alert = (type, msg) => {
            const alert = document.createElement("p");

            alert.classList.add("alert");
            alert.classList.add(`alert-${type}`);
            alert.innerHTML = `<span>${msg}</span>`

            alertContainer.appendChild(alert);

            setTimeout(() => {alert.remove()}, 5000);
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

        const source = new EventSource(window.baseUrl + "/messages");
        source.onopen = (event) => {
            console.debug("/messages sse connection open");
        }
        source.onerror = (event) => {
            console.debug("/messages sse connection error");
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

    }
)

