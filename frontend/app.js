document.getElementById('addCardBtn').addEventListener('click', () => {
    const cardName = prompt('Enter card name:');
    if (cardName) {
        addCard(cardName);
    }
});

function addCard(name) {
    const card = { name };
    fetch('/api/cards/add', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(card),
    }).then(response => {
        if (response.status === 201) {
            loadCards();
        }
    });
}

function loadCards() {
    fetch('/api/cards')
        .then(response => response.json())
        .then(cards => {
            const container = document.getElementById('cardsContainer');
            container.innerHTML = '';
            cards.forEach(card => {
                const div = document.createElement('div');
                div.className = 'card';
                div.textContent = card.name;
                container.appendChild(div);
            });
        });
}

// Load cards when the page loads
loadCards();
