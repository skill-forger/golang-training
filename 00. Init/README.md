# Golang Initiation

## Installation

### Windows
1. Download the latest Go installer from the [official website](https://golang.org/dl/)
2. Run the installer and follow the installation wizard
3. By default, Go will be installed in `C:\Go`
4. The installer will automatically add Go to your PATH environment variable

### macOS
1. Download the latest macOS package from the [official website](https://golang.org/dl/)
2. Open the package and follow the installation prompts
3. Go will be installed in `/usr/local/go`
4. To add Go to your PATH, add the following line to your `~/.bash_profile` or `~/.zshrc`:
   ```
   export PATH=$PATH:/usr/local/go/bin
   ```

### Linux
1. Download the Linux tarball from the [official website](https://golang.org/dl/)
2. Extract it to `/usr/local`:
   ```
   sudo tar -C /usr/local -xzf go1.x.x.linux-amd64.tar.gz
   ```
   (Replace `x.x` with the version number)
3. Add Go to your PATH by adding this line to your `~/.profile` or `~/.bashrc`:
   ```
   export PATH=$PATH:/usr/local/go/bin
   ```

### Verifying Installation
To verify your Go installation, open a terminal or command prompt and run:
   ```bash
   go version
   ```

## Code Editor

### VS Code Setup (Recommended for beginners):
1. Install [Visual Studio Code](https://code.visualstudio.com/)
2. Install the "Go" extension by Google
3. When prompted, install all the recommended Go tools

### Goland JetBrains (Paid IDE)
1. Install [Goland](https://www.jetbrains.com/go/)
2. Sign Up and activate JetBrains account