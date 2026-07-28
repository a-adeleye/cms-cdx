/* =========================================================
   Blog behaviour — self-contained.
   These pages do not load the marketing script, so the shell
   behaviour (nav, reveals, frosted header, subscribe form)
   lives here alongside the blog-specific pieces.
   ========================================================= */

/* ---- Mobile navigation --------------------------------- */
const nav = document.querySelector(".nav");
const navToggle = document.querySelector(".nav-toggle");

navToggle?.addEventListener("click", () => {
  const open = nav.classList.toggle("open");
  navToggle.setAttribute("aria-expanded", String(open));
});

/* ---- Scroll reveal ------------------------------------- */
// Tag the children of any [data-stagger] group so they cascade into view.
document.querySelectorAll("[data-stagger]").forEach((group) => {
  Array.from(group.children).forEach((child, i) => {
    child.classList.add("reveal");
    child.style.setProperty("--reveal-delay", `${i * 90}ms`);
  });
});

const revealObserver = new IntersectionObserver((entries) => {
  entries.forEach((entry) => {
    if (!entry.isIntersecting) return;
    entry.target.classList.add("visible");
    revealObserver.unobserve(entry.target);
  });
}, { threshold: 0.12, rootMargin: "0px 0px -6% 0px" });

document.querySelectorAll(".reveal").forEach((el) => revealObserver.observe(el));

/* ---- Frosted header on scroll -------------------------- */
// A sentinel + IntersectionObserver rather than window.scrollY, because
// `overflow-x: hidden` on <html> makes it the scroll container, which leaves
// window scroll tracking unreliable across browsers.
const siteHeader = document.querySelector(".site-header");

if (siteHeader) {
  const sentinel = document.createElement("div");
  sentinel.setAttribute("aria-hidden", "true");
  sentinel.style.cssText = "position:absolute;top:0;left:0;width:1px;height:12px;pointer-events:none;";
  document.body.prepend(sentinel);
  new IntersectionObserver(([entry]) => {
    siteHeader.classList.toggle("scrolled", !entry.isIntersecting);
  }, { threshold: 0 }).observe(sentinel);
}

/* ---- Subscribe forms ----------------------------------- */
document.querySelectorAll(".subscribe-form").forEach((form) => {
  form.addEventListener("submit", (event) => {
    event.preventDefault();
    const input = event.currentTarget.querySelector("input");
    if (!input.value.trim()) return;
    event.currentTarget.innerHTML = "<strong style='padding:14px;display:block'>You're on the list.</strong>";
  });
});

/* ---- Category filter (all articles page) --------------- */
const chips = document.querySelectorAll(".chip[data-filter]");
const articleGrid = document.querySelector("#articleGrid");
const emptyState = document.querySelector("#emptyState");

if (chips.length && articleGrid) {
  const cards = Array.from(articleGrid.querySelectorAll("[data-category]"));

  chips.forEach((chip) => {
    chip.addEventListener("click", () => {
      chips.forEach((c) => c.classList.remove("active"));
      chip.classList.add("active");

      const filter = chip.dataset.filter;
      let shown = 0;

      cards.forEach((card) => {
        const match = filter === "all" || card.dataset.category === filter;
        card.hidden = !match;
        if (!match) return;

        // Opt the card out of the scroll-reveal rules entirely — toggling
        // `visible` mid-transition can leave a card stranded at opacity 0 —
        // and run its own cascade instead.
        card.classList.remove("reveal", "filter-in");
        card.classList.add("visible");
        void card.offsetWidth; // restart the animation
        card.style.animationDelay = `${shown * 45}ms`;
        card.classList.add("filter-in");
        shown += 1;
      });

      if (emptyState) emptyState.hidden = shown > 0;
    });
  });
}

/* ---- Reading progress (article page) ------------------- */
// `overflow-x: hidden` on html/body makes <html> the scroll container, so read
// the offset from documentElement first and fall back to window.scrollY.
const progressBar = document.querySelector("#readProgress");

if (progressBar) {
  const updateProgress = () => {
    const doc = document.documentElement;
    const scrolled = doc.scrollTop || window.scrollY || 0;
    const max = doc.scrollHeight - doc.clientHeight;
    const ratio = max > 0 ? Math.min(scrolled / max, 1) : 0;
    progressBar.style.transform = `scaleX(${ratio})`;
  };

  let ticking = false;
  const onScroll = () => {
    if (ticking) return;
    ticking = true;
    requestAnimationFrame(() => {
      updateProgress();
      ticking = false;
    });
  };

  window.addEventListener("scroll", onScroll, { passive: true });
  window.addEventListener("resize", onScroll, { passive: true });
  updateProgress();
}

/* ---- Table of contents highlighting -------------------- */
const tocLinks = document.querySelectorAll("#tocList a");

if (tocLinks.length) {
  const headings = Array.from(tocLinks)
    .map((link) => document.querySelector(link.getAttribute("href")))
    .filter(Boolean);

  const setActive = (id) => {
    tocLinks.forEach((link) => {
      link.classList.toggle("active", link.getAttribute("href") === `#${id}`);
    });
  };

  // The last heading to have crossed the reading line is the one being read.
  // An IntersectionObserver callback only carries headings whose visibility
  // *changed*, so during a fast scroll the topmost entry in a batch can belong
  // to a section nowhere near the viewport and the highlight jumps around.
  let ticking = false;
  const update = () => {
    if (ticking) return;
    ticking = true;
    requestAnimationFrame(() => {
      const line = window.innerHeight * 0.18;
      let current = headings[0];
      for (const heading of headings) {
        if (heading.getBoundingClientRect().top > line) break;
        current = heading;
      }
      if (current) setActive(current.id);
      ticking = false;
    });
  };

  window.addEventListener("scroll", update, { passive: true });
  window.addEventListener("resize", update, { passive: true });
  update();
}

/* ---- Share actions ------------------------------------- */
document.querySelectorAll("[data-share]").forEach((button) => {
  button.addEventListener("click", async () => {
    const url = window.location.href;
    const title = document.title;

    if (button.dataset.share === "copy") {
      if (!(await copyToClipboard(url))) return; // blocked — leave the UI alone
      button.classList.add("copied");
      button.setAttribute("aria-label", "Link copied");
      setTimeout(() => {
        button.classList.remove("copied");
        button.setAttribute("aria-label", "Copy link");
      }, 1800);
      return;
    }

    const targets = {
      x: `https://twitter.com/intent/tweet?text=${encodeURIComponent(title)}&url=${encodeURIComponent(url)}`,
      linkedin: `https://www.linkedin.com/sharing/share-offsite/?url=${encodeURIComponent(url)}`
    };

    const target = targets[button.dataset.share];
    if (target) window.open(target, "_blank", "noopener,width=620,height=560");
  });
});

// The async clipboard rejects on insecure origins and when the document has
// lost focus, so fall back to the old selection-based copy.
async function copyToClipboard(text) {
  try {
    await navigator.clipboard.writeText(text);
    return true;
  } catch { /* fall through */ }

  const field = document.createElement("textarea");
  field.value = text;
  field.setAttribute("readonly", "");
  field.style.cssText = "position:fixed;top:0;left:-9999px;opacity:0";
  document.body.appendChild(field);
  field.select();

  let ok = false;
  try { ok = document.execCommand("copy"); } catch { ok = false; }
  field.remove();
  return ok;
}
