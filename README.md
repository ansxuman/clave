<div align="center">
 <img src="https://raw.githubusercontent.com/ansxuman/clave/refs/heads/main/build/appicon.png" alt="Clave Logo" width="120" />
  <h1>Clave</h1>
  <p>A SIMPLE, SECURE, and FREE lightweight cross-platform desktop authenticator app packed with powerful features.</p>

[![License](https://img.shields.io/github/license/ansxuman/clave)](https://github.com/ansxuman/clave/blob/main/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/ansxuman/clave)](https://github.com/ansxuman/clave/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/ansxuman/clave)](https://github.com/ansxuman/clave/issues)
[![GitHub release](https://img.shields.io/github/v/release/ansxuman/clave)](https://github.com/ansxuman/clave/releases)
[![Build Release](https://github.com/ansxuman/clave/actions/workflows/build-release.yml/badge.svg)](https://github.com/ansxuman/clave/actions/workflows/build-release.yml)
[![Notarized by Apple](https://img.shields.io/badge/Release_Notarized_by_Apple-000000?style=flat-square&logo=apple&logoColor=white)](https://developer.apple.com/documentation/security/notarizing-macos-software-before-distribution)

<a href="https://www.producthunt.com/posts/clave-3?embed=true&utm_source=badge-featured&utm_medium=badge&utm_souce=badge-clave&#0045;3" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=669391&theme=light" alt="Clave - Simple&#0044;&#0032;secure&#0032;authentication | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>

![Clave](https://github.com/user-attachments/assets/ac80de84-77a3-48af-ab15-e91afb8a7664)


</div>

## **Features**

- **Easy Profile Addition**:  
  Add accounts via manual entry or QR code image import.

- **Multi-Factor Authentication (MFA)**:  
  Protects access to your TOTP codes with an additional layer of security using a **PIN** or **Touch ID** (macOS only).

- **User-Friendly Interface**:  
  Designed with simplicity in mind, ensuring an intuitive and hassle-free user experience.

- **System Tray Integration**:  
  Quickly access your authenticator from the system tray for seamless usability.

- **Cross-Platform**:  
  Available for **macOS**, **Windows**, and **Linux** .

- **Import & Export**:  
  Easily backup and restore your profiles.

## Development

### Prerequisites
- **NPM**
- **Go**
- **Task**

### Install Wails

1. Clone the Wails repository:
   ```bash
   git clone https://github.com/ansxuman/wails.git
   ```
2. Navigate to the Wails directory:
   ```bash
   cd wails
   ```
3. Check out the `start_on_login` branch:
    ```bash
    git checkout start_on_login
    ```
4. Go to the wails3 directory:
    ```bash
    cd v3/cmd/wails3
    ```
5. Install Wails:
    ```bash
    go install
    ```
    
### Set Up Clave

1. Clone the Clave repository:
   ```bash
   https://github.com/ansxuman/clave.git
   ```
2. Navigate to the Clave directory:
   ```bash
   cd clave
   ```
3. Update the `go.mod` file in the Clave project to replace the Wails path with the local path of your cloned Wails repository.
4. Run the App in Dev Mode
   ```bash
   task dev
   ```

## Contributing

Contributions are what make the open-source community an incredible space for learning, inspiration, and creativity. Any contribution you make is deeply appreciated.Please see our [contributing guidelines](./.github/CONTRIBUTING.md) for more information.

## **Donations**

If you find Clave helpful, consider supporting its development and future updates. Every contribution helps!  

<a href="https://buymeacoffee.com/ansxuman" target="_blank">
<img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;">
</a>

---

## **License**

This project is open-source and available under the [MIT License](LICENSE).
