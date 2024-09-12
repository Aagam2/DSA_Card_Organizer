document.addEventListener("DOMContentLoaded", function() {
    const urlParams = new URLSearchParams(window.location.search);
    const algorithmId = urlParams.get("algorithmId");
    const algorithmName = urlParams.get("algorithmName");

    const algorithmTitle = document.getElementById("algorithm-title");
    const algorithmCodeElement = document.getElementById("algorithm-code");
    const notesContent = document.getElementById("notes-content");
    const saveNotesButton = document.getElementById("save-notes");
    const resizer = document.getElementById("resizer");

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

    // Function to fetch and display notes
    function fetchAndDisplayNotes() {
        fetch(`/api/notes/${algorithmId}`)
            .then(response => response.json())
            .then(data => {
                if (data.notes) {
                    notesContent.value = data.notes;
                } else {
                    notesContent.value = "";
                }
            })
            .catch(error => {
                console.error("Error fetching notes:", error);
                notesContent.value = "Error fetching notes.";
            });
    }

    // Function to save notes
    function saveNotes() {
        const notes = notesContent.value;
        fetch(`/api/notes/${algorithmId}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ notes })
        })
        .then(response => response.json())
        .then(data => {
            alert("Notes saved successfully!");
        })
        .catch(error => {
            console.error("Error saving notes:", error);
            alert("Error saving notes.");
        });
    }

    // Fetch and display the algorithm code and notes on page load
    fetchAndDisplayAlgorithmCode();
    fetchAndDisplayNotes();

    // Save notes on button click
    saveNotesButton.addEventListener("click", saveNotes);

    // Resizing logic
    const container = document.getElementById("container");
    let isResizing = false;

    resizer.addEventListener('mousedown', (e) => {
        isResizing = true;
        document.addEventListener('mousemove', handleMouseMove);
        document.addEventListener('mouseup', () => {
            isResizing = false;
            document.removeEventListener('mousemove', handleMouseMove);
        });
    });

    function handleMouseMove(e) {
        if (!isResizing) return;
        const containerRect = container.getBoundingClientRect();
        const totalWidth = containerRect.width;
        const minAlgorithmWidth = totalWidth * 0.5; // Minimum 50% for algorithm section
        const maxAlgorithmWidth = totalWidth * 0.85; // Maximum 85% for algorithm section
    
        let newAlgorithmWidth = e.clientX - containerRect.left;
        newAlgorithmWidth = Math.max(minAlgorithmWidth, Math.min(newAlgorithmWidth, maxAlgorithmWidth));
        
        const notesWidth = totalWidth - newAlgorithmWidth - 10; // 10px for resizer
        
        container.style.gridTemplateColumns = `${newAlgorithmWidth}px 10px ${notesWidth}px`;
    }   
});
