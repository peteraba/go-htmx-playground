// movies
$(function () {
    const filmsCheckAll = $("#check-all");
    if (!filmsCheckAll) {
        console.warn("check-all element not found.");
        return;
    }

    filmsCheckAll.parent().removeClass("hidden");

    const filmChecks = $("#film-list .film-check");
    if (!filmChecks) {
        console.debug("no films found");
        return;
    }

    filmsCheckAll.on("click", (e) => {
        const checked = e.currentTarget.checked;
        console.info(`checked: ${checked}`);
        filmChecks.each((idx) => {
            $(filmChecks[idx]).attr("checked", checked);
        });
    })
})
