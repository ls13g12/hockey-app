import { mockedMsg } from "../mocked_data/mocked_message";

//async/await type function
export async function GetWelcomeMsg() {
  try {
    const resp = await fetch('http://localhost:3333/hello', {
      "Content-Type": "application/json"
    });
    if (resp.ok) {
      const data = await resp.json();
      return data.welcome_msg;
    } else {
      return mockedMsg;
    }
  } catch (err) {
    console.error(`error when fetching welcome message: ${err}`);
    return mockedMsg;
  }
}
