document.getElementById("loginForm").addEventListener("submit", function(event) {
    event.preventDefault();
  
    const email = document.getElementById("exampleInputEmail1").value;
    const password = document.getElementById("exampleInputPassword1").value;
  
    const data = {
      email: email,
      password: password
    };
  
    fetch("/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
      if (data.token) {  // Передбачається, що сервер повертає токен у полі `token`
        window.location.reload();  // Оновлюємо сторінку, щоб відобразити кнопку "Upload"
      } else {
        console.log("Authentication failed");
      }
    })
    .catch((error) => {
      console.error("Помилка:", error);
    });
  });
  