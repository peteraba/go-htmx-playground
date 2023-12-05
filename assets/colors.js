// set page color
$(function () {
    const html = $("#html");
    if (!html) {
        console.warn("no #html element found.");
        return;
    }

    const container = $("#color-selectors");
    if (!container) {
        console.debug("no #color-selectors element found.");
        return;
    }

    const colorSelectors = $(".color-selector");
    if (colorSelectors.length === 0) {
        console.debug("no .color-selector elements found.");
        return;
    }

    colorSelectors.on("click", (e) => {
        console.log(e)
        const target = $(e.currentTarget);
        const theme = target.data("theme");
        if (theme) {
            html.attr("data-theme", theme);
            localStorage.setItem("theme", theme);
        }
    })
})