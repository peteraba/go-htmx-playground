// set page color
ready(() => {
    const html = document.getElementById("html");
    if (!html) {
        console.warn("no #html element found.");
        return;
    }

    const container = document.getElementById("color-selectors");
    if (!container) {
        console.debug("no #color-selectors element found.");
        return;
    }

    const colorSelectors = container.querySelectorAll(".color-selector");
    if (!colorSelectors) {
        console.debug("no .color-selector elements found.");
        return;
    }

    colorSelectors.forEach(cs => {
        cs.addEventListener("click", (e) => {
            const theme = cs.getAttribute("data-theme")
            html.setAttribute("data-theme", theme);
            localStorage.setItem("theme", theme);
        })
    })

    console.info("listener is set up on " + colorSelectors.length + " color selectors.");
});