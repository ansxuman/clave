const themeToggle = document.getElementById('themeToggle');
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)');

function setTheme(isDark) {
  document.documentElement.setAttribute('data-theme', isDark ? 'dark' : 'light');
  const themeIcon = themeToggle?.querySelector('.theme-icon');
  if (themeIcon) {
    themeIcon.textContent = isDark ? '☀️' : '🌙';
  }
  localStorage.setItem('theme', isDark ? 'dark' : 'light');
}

const savedTheme = localStorage.getItem('theme');
if (savedTheme) {
  setTheme(savedTheme === 'dark');
} else {
  setTheme(prefersDark.matches);
}

themeToggle?.addEventListener('click', () => {
  const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
  setTheme(!isDark);
});

prefersDark.addEventListener('change', (e) => {
  if (!localStorage.getItem('theme')) {
    setTheme(e.matches);
  }
});

let currentVersion = '';

async function updateVersion() {
  try {
    const response = await fetch('https://api.github.com/repos/ansxuman/clave/releases/latest');
    if (!response.ok) {
      throw new Error('Failed to fetch releases');
    }
    const data = await response.json();
    if (!data || !data.tag_name) {
      throw new Error('Invalid release data');
    }
    
    currentVersion = data.tag_name.replace(/^v/, '');
    if (currentVersion) {
      updateDownloadLinks(currentVersion);
    }
  } catch (error) {
    console.error('Failed to fetch version:', error);
  }
}

function updateDownloadLinks(versionNumber) {
  const platformUrls = {
    'mac-universal': `Clave-v${versionNumber}-universal.dmg`,
    'windows': `Clave-Setup-v${versionNumber}-x64.exe`,
    'linuxAppImage': `clave_v${versionNumber}_amd64.AppImage`,
    'linuxDeb': `clave_v${versionNumber}_amd64.deb`,
  };

  document.querySelectorAll('.download-button').forEach(link => {
    const platform = link.getAttribute('data-platform');
    if (platform && platformUrls[platform]) {
      link.href = `https://github.com/ansxuman/clave/releases/download/v${versionNumber}/${platformUrls[platform]}`;
    }
  });
}

document.querySelectorAll('.download-button').forEach(button => {
  button.addEventListener('click', (e) => {
    if (typeof gtag === 'function') {
      gtag('event', 'download', {
        'event_category': 'App',
        'event_label': button.getAttribute('data-platform'),
        'value': currentVersion
      });
    }
  });
});

document.querySelectorAll('nav a[href^="#"]').forEach(anchor => {
  anchor.addEventListener('click', function (e) {
    e.preventDefault();
    const target = document.querySelector(this.getAttribute('href'));
    if (target) {
      target.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  });
});

document.addEventListener('DOMContentLoaded', () => {
  if (typeof AOS !== 'undefined') {
    AOS.init({
      duration: 800,
      once: true,
      offset: 100
    });
  }
  updateVersion();
});

const smoothScroll = (target) => {
  const element = document.querySelector(target);
  element?.scrollIntoView({ 
    behavior: 'smooth',
    block: 'start'
  });
};