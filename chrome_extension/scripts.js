// chrome.browserAction.onClicked.addListener(function(tab) {
//   chrome.tabs.executeScript(tab.id,{code:   `document.getElementById("productTitle").innerText`},sendCurrentTitle);
//  });
function go(url) {
    var html = document.documentElement.innerHTML
    alert(url)
    const params = {
        content: html,
        website: url
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