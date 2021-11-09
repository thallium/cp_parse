function go(url) {
    var html = document.documentElement.innerHTML
    const params = {
        body: html,
        url: url
    }
    const options = {
        method: 'POST',
        body: JSON.stringify(params)
    }
    fetch('http://localhost:8090', options)
}
chrome.action.onClicked.addListener(async (tab) => {
    url = tab.url

    chrome.scripting.executeScript({
        target: { tabId: tab.id },
        func: go,
        args: [url],
    });
});