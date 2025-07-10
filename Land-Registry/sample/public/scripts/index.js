const registerLand = async (event) => {
    event.preventDefault();

    const landId = document.getElementById("landId").value;
    const location = document.getElementById("location").value;
    const size = document.getElementById("size").value;
    const owner = document.getElementById("owner").value;

    const landData = {
        landId,
        location,
        size,
        owner,
    };

    if (!landId || !location || !size || !owner) {
        alert("Please enter all land details.");
        return;
    }

    try {
        const response = await fetch("/api/land", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(landData),
        });

        const data = await response.json();
        console.log("Land Registration Response:", data);

        if (response.ok) {
            alert(" Land registered successfully.");
        } else {
            alert(" Failed to register land: " + (data.message || "Unknown error."));
        }
    } catch (error) {
        alert(" Error registering land.");
        console.error(error);
    }
};

const getLand = async (event) => {
    if (event) event.preventDefault();

    const landId = document.getElementById("landIdInput").value;
    if (!landId) {
        alert("Please enter a valid Land ID.");
        return;
    }

    try {
        const response = await fetch(`/api/land/${landId}`);
        const text = await response.text();

        if (!response.ok) {
            alert(" Failed to fetch land info.");
            console.error("Error:", text);
            document.getElementById("output").textContent = `Not found: ${landId}`;
            return;
        }

        const data = JSON.parse(text);
        console.log("Land Data:", data);
        document.getElementById("output").textContent = JSON.stringify(data.data, null, 2);
    } catch (error) {
        alert("Error fetching land info.");
        console.error(error);
    }
};

document.addEventListener("DOMContentLoaded", () => {
    const registerForm = document.getElementById("registerForm");
    const getForm = document.getElementById("getForm");

    if (registerForm) registerForm.addEventListener("submit", registerLand);
    if (getForm) getForm.addEventListener("submit", getLand);
});
