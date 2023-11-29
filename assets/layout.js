const ready = (callback) => {
    if (document.readyState !== "loading") callback();
    else document.addEventListener("DOMContentLoaded", callback);
}

// set page color
ready(() => {
    const html = document.getElementById("html");
    if (!html) {
        console.warn("no #html element found.");
        return;
    }

    const theme = localStorage.getItem("theme");
    if (theme) {
        html.setAttribute("data-theme", theme);
        console.debug("custom theme is set: " + theme)
    }
});