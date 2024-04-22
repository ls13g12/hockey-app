//async/await type function
export async function GetBackendHealth() {
  try {
    const resp = await fetch('http://localhost:8080/healthcheck', {
      "Content-Type": "application/json"
    });
    if (resp.ok) {
      const data = await resp.json();
      return `healthcheck is ${data.message}`;
    } else {
      return `error reaching backend`;
    }
  } catch (err) {
    return `error reaching backend`;
  }
}
