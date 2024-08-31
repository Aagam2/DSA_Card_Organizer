document.addEventListener("DOMContentLoaded", function() {
    const cardsContainer = document.getElementById("cards-container");
    const addCardForm = document.getElementById("add-card-form");
    const cardNameInput = document.getElementById("card-name");

    // Function to fetch and display cards
    function fetchAndDisplayCards() {
        fetch("/api/cards")
            .then(response => response.json())
            .then(data => {
                cardsContainer.innerHTML = ""; // Clear existing cards
                data.forEach(card => {
                    const cardElement = document.createElement("div");
                    cardElement.className = "card";
                    cardElement.innerHTML = `
                        <h3>${card.name}</h3>
                        <p>ID: ${card.id}</p>
                    `;
                    cardElement.addEventListener("click", () => {
                        window.location.href = `subtopics.html?topicId=${card.id}&topicName=${card.name}`;
                    });
                    cardsContainer.appendChild(cardElement);
                });
            })
            .catch(error => {
                console.error("Error fetching cards:", error);
            });
    }

    // Fetch and display cards on page load
    fetchAndDisplayCards();

    // Handle form submission to add new card
    addCardForm.addEventListener("submit", function(event) {
        event.preventDefault();
        const cardName = cardNameInput.value;

        fetch("/api/cards/add", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ name: cardName })
        })
        .then(response => response.json())
        .then(data => {
            cardNameInput.value = ""; // Clear the input field
            fetchAndDisplayCards(); // Refresh the list of cards
        })
        .catch(error => {
            console.error("Error adding card:", error);
        });
    });
});