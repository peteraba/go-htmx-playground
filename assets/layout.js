// Theme switcher
$(function () {
    const theme = localStorage.getItem("theme");

    $("#html").attr("data-theme", theme);
});

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

// Navigation changes
$(function () {
    const listItems = $("#topnav > li > a");

    const onChangeState = (state, title, url) => {
        listItems.each((idx, listItem) => {
            const $li = $(listItem);

            if ($li.attr('hx-get') === url) {
                $li.addClass("link-primary");
            } else {
                $li.removeClass("link-primary");
            }
        })
    }

    // set onChangeState() listener:
    ['pushState'].forEach((changeState) => {
        // store original values under underscored keys (`window.history._pushState()` and `window.history._replaceState()`):
        window.history['_' + changeState] = window.history[changeState]

        window.history[changeState] = new Proxy(window.history[changeState], {
            apply (target, thisArg, argList) {
                const [state, title, url] = argList
                onChangeState(state, title, url)

                return target.apply(thisArg, argList)
            },
        })
    })
});