document.addEventListener("DOMContentLoaded", function() {
    const urlParams = new URLSearchParams(window.location.search);
    const topicId = urlParams.get("topicId");
    const topicName = urlParams.get("topicName");

    const topicTitle = document.getElementById("topic-title");
    const subtopicsContainer = document.getElementById("subtopics-container");
    const addSubtopicForm = document.getElementById("add-subtopic-form");
    const subtopicNameInput = document.getElementById("subtopic-name");
    const subtopicDescriptionInput = document.getElementById("subtopic-description");

    topicTitle.textContent = `Subtopics for ${topicName}`;

    // Function to fetch and display subtopics
    async function fetchAndDisplaySubtopics() {
        try {
            const response = await fetch(`/api/subtopics/${topicId}`);
            const data = await response.json();
            subtopicsContainer.innerHTML = ""; // Clear existing subtopics
            data.forEach(subtopic => {
                const subtopicElement = document.createElement("div");
                subtopicElement.className = "subtopic card";
                subtopicElement.innerHTML = `
                    <h3>${subtopic.name}</h3>
                    <p>${subtopic.description}</p>
                `;
                subtopicElement.addEventListener("click", function() {
                    updateNavigation({ subtopicId: subtopic.id, subtopicName: subtopic.name });
                    window.location.href = `algorithms.html?subtopicId=${subtopic.id}&subtopicName=${subtopic.name}`;
                });
                subtopicsContainer.appendChild(subtopicElement);
            });
        } catch (error) {
            console.error("Error fetching subtopics:", error);
        }
    }

    // Fetch and display subtopics on page load
    fetchAndDisplaySubtopics();

    // Handle form submission to add new subtopic
    addSubtopicForm.addEventListener("submit", async function(event) {
        event.preventDefault();
        const subtopicName = subtopicNameInput.value;
        const subtopicDescription = subtopicDescriptionInput.value;

        try {
            const response = await fetch("/api/subtopics/add", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({ topicId: topicId, name: subtopicName, description: subtopicDescription })
            });
            const newSubtopic = await response.json();
            subtopicNameInput.value = ""; // Clear the input fields
            subtopicDescriptionInput.value = "";

            // Add the new subtopic to the UI
            const subtopicElement = document.createElement("div");
            subtopicElement.className = "subtopic card";
            subtopicElement.innerHTML = `
                <h3>${newSubtopic.name}</h3>
                <p>${newSubtopic.description}</p>
            `;
            subtopicElement.addEventListener("click", function() {
                window.location.href = `algorithms.html?subtopicId=${newSubtopic.id}&subtopicName=${newSubtopic.name}`;
            });
            subtopicsContainer.appendChild(subtopicElement);
        } catch (error) {
            console.error("Error adding subtopic:", error);
        }
    });
});