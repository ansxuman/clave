# Clave

Clave is a secure desktop authenticator app that allows you to easily manage and generate time-based one-time passwords (TOTP) for your accounts, with an additional layer of security.

## Features

- **Three Easy Ways to Add Profiles**: Manual entry, drag & drop QR code, or import QR code image.
- **Secure Authentication**: Generates time-based one-time passwords (TOTP) for your accounts.
- **Multi-Factor Authentication**: Requires PIN or Touch ID (macOS only) to access TOTP codes.
- **User-Friendly Interface**: Simple and intuitive design for easy navigation.
- **System Tray Integration**: Easily accessible from the system tray for quick use.
- **Cross-Platform**: Available for macOS, Windows, and Linux.

## Usage

1. Launch Clave. The app will appear in your system tray.
2. Click on the Clave icon in the system tray to open the main window.
3. Choose your preferred method to add a new profile:
   - Manual entry: Input your secret key and issuer details.
   - Drag & drop QR: Drag and drop a QR code image.
   - Import QR code: Select a QR code image from your device.
4. To access your generated TOTP codes:
   - Open the app from the system tray.
   - Enter your PIN or use Touch ID (macOS only) for authentication.
   
## Technology Stack

- [Wails](https://wails.io/): Application framework
- [Go](https://golang.org/): Client-side Backend language
- [Svelte](https://svelte.dev/): Frontend framework
- [Tailwind CSS](https://tailwindcss.com/): Styling


## License

This project is licensed under the [MIT License](LICENSE.txt).

## Acknowledgements

- [Wails](https://wails.io/) for the application framework.
- [Inter](https://github.com/rsms/inter) for the font used in the application.