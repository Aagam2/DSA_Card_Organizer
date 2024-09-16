function storeNavData(key, value) {
    localStorage.setItem(key, value);
}

function getNavData(key) {
    return localStorage.getItem(key);
}

function createNavbar() {
    const navbar = document.createElement('nav');
    const navItems = [
        { text: 'Home', link: 'index.html' },
    ];

    const topicId = getNavData('topicId');
    const topicName = getNavData('topicName');
    const subtopicId = getNavData('subtopicId');
    const subtopicName = getNavData('subtopicName');
    const algorithmId = getNavData('algorithmId');
    const algorithmName = getNavData('algorithmName');

    // if (topicId && topicName) {
    //     navItems.push({ text: topicName, link: `subtopics.html?topicId=${topicId}&topicName=${topicName}` });
    // }

    if (subtopicId && subtopicName) {
        // navItems.push({ text: subtopicName, link: `algorithms.html?subtopicId=${subtopicId}&subtopicName=${subtopicName}` });
        navItems.push({ text: topicName, link: `subtopics.html?topicId=${topicId}&topicName=${topicName}` });
    }

    if (algorithmId && algorithmName) {
        // navItems.push({ text: algorithmName, link: `algorithm.html?algorithmId=${algorithmId}&algorithmName=${algorithmName}` });
        navItems.push({ text: subtopicName, link: `algorithms.html?subtopicId=${subtopicId}&subtopicName=${subtopicName}` });
    }

    navbar.innerHTML = `
        <ul>
            ${navItems.map(item => `<li><a href="${item.link}">${item.text}</a></li>`).join('')}
        </ul>
    `;
    document.body.insertBefore(navbar, document.body.firstChild);
}

document.addEventListener('DOMContentLoaded', createNavbar);

// Call this function when navigating to a new page
function updateNavigation(params) {
    for (const [key, value] of Object.entries(params)) {
        console.log(key, value);
        storeNavData(key, value);
    }
}

// Example usage:
// updateNavigation({ topicId: '123', topicName: 'Data Structures' });
