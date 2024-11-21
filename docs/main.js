const themeToggle = document.getElementById('themeToggle');
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)');

function setTheme(isDark) {
  document.documentElement.setAttribute('data-theme', isDark ? 'dark' : 'light');
  themeToggle.querySelector('.theme-icon').textContent = isDark ? '☀️' : '🌙';
  localStorage.setItem('theme', isDark ? 'dark' : 'light');
}

const savedTheme = localStorage.getItem('theme');
if (savedTheme) {
  setTheme(savedTheme === 'dark');
} else {
  setTheme(prefersDark.matches);
}

themeToggle.addEventListener('click', () => {
  const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
  setTheme(!isDark);
});

prefersDark.addEventListener('change', (e) => {
  if (!localStorage.getItem('theme')) {
    setTheme(e.matches);
  }
});

async function updateVersion() {
  try {
    const response = await fetch('https://api.github.com/repos/ansxuman/clave/releases/latest');
    const data = await response.json();
    const version = data.tag_name;
    const versionNumber = version.match(/\d+\.\d+\.\d+/)?.[0];

    if (versionNumber) {
      document.getElementById('current-version').textContent = "v" + versionNumber;
      updateDownloadLinks(versionNumber);
    }
  } catch (error) {
    console.error('Failed to fetch version:', error);
  }
}

function updateDownloadLinks(versionNumber) {
  const platformUrls = {
    'macos-intel': `intel-Clave_${versionNumber}_x64.dmg`,
    'macos-silicon': `silicon-Clave_${versionNumber}_aarch64.dmg`,
    'windows': `Clave_${versionNumber}_x64.exe`,
    'linux-deb-amd64': `Clave_${versionNumber}_amd64.deb`,
  };

  document.querySelectorAll('.download-button').forEach(link => {
    const platform = link.getAttribute('data-type');
    if (platformUrls[platform]) {
      link.href = `https://github.com/ansxuman/clave/releases/download/v${versionNumber}/${platformUrls[platform]}`;
    }
  });
}

document.querySelectorAll('.faq-question').forEach(button => {
  button.addEventListener('click', () => {
    const faqItem = button.parentElement;
    const isActive = faqItem.classList.contains('active');
    document.querySelectorAll('.faq-item').forEach(item => item.classList.remove('active'));
    if (!isActive) faqItem.classList.add('active');
  });
});

document.querySelectorAll('.download-button').forEach(button => {
  button.addEventListener('click', (e) => {
    gtag('event', 'download', {
      'event_category': 'App',
      'event_label': button.getAttribute('data-type'),
      'value': button.getAttribute('data-version')
    });
  });
});

document.querySelectorAll('a[href^="#"]').forEach(anchor => {
  anchor.addEventListener('click', function (e) {
    e.preventDefault();
    const target = document.querySelector(this.getAttribute('href'));
    if (target) {
      target.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  });
});

const observer = new IntersectionObserver(
  (entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        entry.target.classList.add('visible');
      }
    });
  },
  { threshold: 0.1 }
);

document.querySelectorAll('.feature, .step, .download-button').forEach(el => {
  observer.observe(el);
});


document.addEventListener('DOMContentLoaded', () => {
  updateVersion();
  fetchDownloadStats();
});



document.addEventListener('DOMContentLoaded', () => {
  const menuButton = document.querySelector('.menu-button');
  const navLinks = document.querySelector('.nav-links');

  menuButton.addEventListener('click', () => {
    navLinks.classList.toggle('active');
  });

  document.addEventListener('click', (e) => {
    if (!navLinks.contains(e.target) && !menuButton.contains(e.target)) {
      navLinks.classList.remove('active');
    }
  });
});

AOS.init({
  duration: 800,
  once: true,
  offset: 100
});

const smoothScroll = (target) => {
  const element = document.querySelector(target);
  element?.scrollIntoView({ 
    behavior: 'smooth',
    block: 'start'
  });
};