const preview = document.getElementById("previewImage");
const device = document.getElementById("device");
const lang = document.getElementById("lang");
const tz = document.getElementById("timezone");
const weekends = document.getElementById("weekends");
const urlBox = document.getElementById("url");
const copyBtn = document.getElementById("copy");

function buildURL() {
    const params = new URLSearchParams({
        device: device.value,
        lang: lang.value,
        timezone: tz.value,
        weekends: weekends.value,
    });
    return `/wallpaper?${params.toString()}`;
}

function update() {
    const url = buildURL();
    preview.src = url;
    urlBox.textContent = location.origin + url;
}

[device, lang, tz, weekends].forEach(el => {
    el.addEventListener("change", update);
    el.addEventListener("input", update);
});

copyBtn.onclick = () => {
    navigator.clipboard.writeText(urlBox.textContent);
    copyBtn.textContent = "Copied ✓";
    setTimeout(() => (copyBtn.textContent = "Copy link"), 1200);
};

update();
