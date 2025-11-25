document.addEventListener("DOMContentLoaded", async function() {
  const navContainer = document.querySelector(".navbar-nav");

  try {
    const res = await fetch("/admin/check");

    if (res.status === 200) {
      // Token still valid → show Upload
      navContainer.innerHTML = `
        <li class="nav-item">
          <a class="nav-link active" aria-current="page" href="/">Search</a>
        </li>
        <li class="nav-item">
          <a class="nav-link" href="/admin/upload">Upload</a>
        </li>
      `;
    } else {
      throw new Error("unauthorized");
    }

  } catch {
    // Token expired → show Login
    navContainer.innerHTML = `
      <li class="nav-item">
        <a class="nav-link active" aria-current="page" href="/">Search</a>
      </li>
      <li class="nav-item">
        <a class="nav-link" href="login_page.html">Login</a>
      </li>
    `;

    // cleaning, just in case
    localStorage.removeItem("jwt");
  }
});
