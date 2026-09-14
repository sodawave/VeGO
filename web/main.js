document.querySelectorAll('a[href^="#"]').forEach((a) => {
  a.addEventListener("click", (e) => {
    const id = a.getAttribute("href");
    const el = document.querySelector(id);
    if (!el) return;
    e.preventDefault();
    el.scrollIntoView({ behavior: "smooth", block: "start" });
  });
});

const sample = document.querySelector(".hero .sample code");
if (sample) {
  const full = sample.textContent;
  sample.textContent = "";
  let i = 0;
  const tick = () => {
    sample.textContent = full.slice(0, i);
    i += 1;
    if (i <= full.length) requestAnimationFrame(() => setTimeout(tick, 12));
  };
  setTimeout(tick, 400);
}

const panels = document.querySelectorAll(".panel");
if ("IntersectionObserver" in window) {
  const io = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) entry.target.classList.add("in-view");
      });
    },
    { threshold: 0.12 }
  );
  panels.forEach((p) => io.observe(p));
} else {
  panels.forEach((p) => p.classList.add("in-view"));
}
