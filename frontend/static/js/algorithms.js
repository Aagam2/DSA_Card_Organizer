document.addEventListener("DOMContentLoaded", function() {
    const urlParams = new URLSearchParams(window.location.search);
    const subtopicId = urlParams.get("subtopicId");
    const subtopicName = urlParams.get("subtopicName");

    const subtopicTitle = document.getElementById("subtopic-title");
    const algorithmsContainer = document.getElementById("algorithms-container");
    const addAlgorithmForm = document.getElementById("add-algorithm-form");
    const algorithmNameInput = document.getElementById("algorithm-name");
    const algorithmDescriptionInput = document.getElementById("algorithm-description");
    const algorithmCodeInput = document.getElementById("algorithm-code");

    subtopicTitle.textContent = `Algorithms for ${subtopicName}`;

    // Function to fetch and display algorithms
    function fetchAndDisplayAlgorithms() {
        fetch(`/api/algorithms?subtopicId=${subtopicId}`)
            .then(response => response.json())
            .then(data => {
                algorithmsContainer.innerHTML = ""; // Clear existing algorithms
                data.forEach(algorithm => {
                    const algorithmElement = document.createElement("div");
                    algorithmElement.className = "algorithm card";
                    algorithmElement.innerHTML = `
                        <h3>${algorithm.name}</h3>
                        <p>${algorithm.description}</p>
                    `;
                    algorithmElement.addEventListener("click", function() {
                        window.location.href = `algorithm.html?algorithmId=${algorithm.id}&algorithmName=${algorithm.name}`;
                    });
                    algorithmsContainer.appendChild(algorithmElement);
                });
            })
            .catch(error => {
                console.error("Error fetching algorithms:", error);
            });
    }

    // Fetch and display algorithms on page load
    fetchAndDisplayAlgorithms();

    // Handle form submission to add new algorithm
    addAlgorithmForm.addEventListener("submit", function(event) {
        event.preventDefault();
        const algorithmName = algorithmNameInput.value;
        const algorithmDescription = algorithmDescriptionInput.value;
        const algorithmCode = algorithmCodeInput.files[0];

        const formData = new FormData();
        formData.append("subtopicId", subtopicId);
        formData.append("name", algorithmName);
        formData.append("description", algorithmDescription);
        formData.append("code", algorithmCode);

        fetch("/api/algorithms/add", {
            method: "POST",
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            algorithmNameInput.value = ""; // Clear the input fields
            algorithmDescriptionInput.value = "";
            algorithmCodeInput.value = "";
            fetchAndDisplayAlgorithms(); // Refresh the list of algorithms
        })
        .catch(error => {
            console.error("Error adding algorithm:", error);
        });
    });
});