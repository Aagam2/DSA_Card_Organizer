document.addEventListener("DOMContentLoaded", function() {
    const urlParams = new URLSearchParams(window.location.search);
    const algorithmId = urlParams.get("algorithmId");
    const algorithmName = urlParams.get("algorithmName");

    const algorithmTitle = document.getElementById("algorithm-title");
    const algorithmCodeElement = document.getElementById("algorithm-code");

    algorithmTitle.textContent = `Code for ${algorithmName}`;

    // Function to fetch and display the algorithm code
    function fetchAndDisplayAlgorithmCode() {
        fetch(`/api/algorithms/${algorithmId}`)
            .then(response => response.json())
            .then(data => {
                if (data.code) {
                    algorithmCodeElement.textContent = data.code;
                } else {
                    algorithmCodeElement.textContent = "No code available for this algorithm.";
                }
            })
            .catch(error => {
                console.error("Error fetching algorithm code:", error);
                algorithmCodeElement.textContent = "Error fetching algorithm code.";
            });
    }

    // Fetch and display the algorithm code on page load
    fetchAndDisplayAlgorithmCode();
});