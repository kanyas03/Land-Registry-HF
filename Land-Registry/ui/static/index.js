document.getElementById("listLandForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    const form = e.target;
    const data = Object.fromEntries(new FormData(form).entries());
  
    const res = await fetch("/api/list-land", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data)
    });
  
    document.getElementById("listLandResult").innerText = await res.text();
    form.reset();
  });
  
  async function getAvailableLands() {
    const res = await fetch("/api/get-available-lands");
    const data = await res.json();
    document.getElementById("availableLands").innerText = JSON.stringify(data, null, 2);
  }
  
  document.getElementById("requestBuyForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    const form = e.target;
    const { offerID, ...rest } = Object.fromEntries(new FormData(form).entries());
  
    const res = await fetch("/api/request-buy", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        offerID,
        buyerRequest: rest
      })
    });
  
    document.getElementById("requestBuyResult").innerText = await res.text();
    form.reset();
  });
  
  document.getElementById("registerBuyerForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    const form = e.target;
    const { landID, ...ownership } = Object.fromEntries(new FormData(form).entries());
  
    const res = await fetch("/api/register-buyer", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        landID,
        buyerOwnership: ownership
      })
    });
  
    document.getElementById("registerBuyerResult").innerText = await res.text();
    form.reset();
  });
  