<div align="center">

# 🚀 Migration from the 3X-UI Panel to PasarGuard

![Migration](./images/migration.svg)

# User and Inbound Panel Manager

[🇬🇧 English](./README_EN.md) | [🇮🇷 فارسی](./README.md)

<p>
A powerful and comprehensive tool for migrating users between 3X-UI and PasarGuard panels
</p>

![Status](https://img.shields.io/badge/Status-Active-brightgreen)
![Go](https://img.shields.io/badge/Go-1.20%2B-blue)
![License](https://img.shields.io/badge/License-MIT-green)

</div>

---

This is a migration tool between 3X-UI and PasarGuard panels.
With this tool, transfer your VPN user accounts while preserving traffic amount and remaining time between 3X-UI and PasarGuard panels.

**🎯 Main Goal: Transfer users from 3X-UI panel to PasarGuard panel**

## ✨ Features

<div>

- ✅ Export (Backup) all users from 3X-UI panel with or without inbound details
- ✅ Export (Backup) all users from PasarGuard panel
- ✅ Export users in a format compatible with both panels
- ✅ Convert, transfer and import users to both 3X-UI and PasarGuard panels
- ✅ Automatic traffic management during transfer
- ✅ Automatic resolution of username conflicts in panels
- ✅ Assign specific Groups to users during transfer to PasarGuard

</div>

## 🚀 How to Use

### Compilation

```bash
cd /root/Panels_Migration
go build -o Panels_Migration ./cmd/main
./Panels_Migration
```

### Install

```bash
curl -sSL https://raw.githubusercontent.com/ArashAfkandeh/Panels-Migration/main/install.sh | sudo bash && /root/Panels_Migration
```

### Main Steps

The program has a menu-driven interface that includes the following sections:

**📋 MAIN MENU:**
- **🔧 PANEL SELECTION** - Select the desired panel (3X-UI or PasarGuard)
- **🚪 APPLICATION CONTROL** - Exit the program

**3X-UI OPERATIONS:**
- **📤 EXPORT OPERATIONS** - Export users and Inbounds
- **📥 IMPORT OPERATIONS** - Import users and Inbounds
- **🔙 NAVIGATION** - Return to main menu

**PASARGUARD OPERATIONS:**
- **📤 EXPORT OPERATIONS** - Export users
- **📥 IMPORT OPERATIONS** - Import users
- **🔙 NAVIGATION** - Return to main menu

## 📸 Program Menu

![Main Menu](./images/Main_Menu.png)

## 📋 Quick Usage Guide

### Transfer Users from 3X-UI to PasarGuard

**Step One - Export from 3X-UI:**

```
1. Run the program: ./Panels_Migration
2. Select option [1] (3X-UI Panel Operations)
3. Select option [2] (Export users only - PasarGuard format)
4. Enter 3X-UI panel credentials:
   - Address: https://panel.example.com:port/path
   - Username: admin
   - Password: password
5. Specify the export file name (e.g.: users_backup.json)
```

**Step Two - Import to PasarGuard:**

```
1. Run the program again: ./Panels_Migration
2. Select option [2] (PasarGuard Panel Operations)
3. Select option [2] (Import users from JSON)
4. Enter PasarGuard panel credentials:
   - Address: https://panel.example.com:port
   - Username: admin
   - Password: password
5. Enter the path to the previous export file: users_backup.json
6. Select desired groups (optional)
```

## ⚠️ Notes

- To transfer users from 3X-UI panel to PasarGuard panel, you must export users `(without inbound details)`.
- **Supported Protocols:** This tool currently only supports transferring users with **VLESS** and **VMess** protocols.
- **Backup:** Always back up your data before transferring users
- **API Access:** Make sure you have access to the panel APIs and that the firewall doesn't block them
- **User Traffic:** During transfer, user traffic is calculated from the `traffic_remaining` value
- **Naming Conflicts:** If a username is duplicated, the program automatically adds a suffix (e.g.: `user_1`)

## 🔍 Verbose Usage

To see operation details and troubleshooting, run the program with the `-v` argument:

```bash
./Panels_Migration -v
```

This mode displays:
- API request details
- UUID mapping
- Steps to resolve username conflicts
- Traffic calculations

## 💖 Donation

<div align="center">

Support our project by donating TRX (TRC-20 network):

![Wallet QR Code](./images/donation_qrcode.png)

```
THzUk99MRsGDgrYK1Nh1YhA3cVuTosmHoA
```

Thank you for your support! ❤️

</div>

---

<div align="center">

**Easy-to-use tool for transferring and managing users**

**Made with ❤️ by ArashAfkandeh**

</div>
